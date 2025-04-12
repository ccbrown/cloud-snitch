import { init, Plugin, RematchDispatch, RematchRootState } from '@rematch/core';
import loadingPlugin, { ExtraModelsFromLoading } from '@rematch/loading';
import { produce } from 'immer';
import Redux from 'redux';
import { TypedUseSelectorHook, useDispatch as useDispatchImpl, useSelector as useSelectorImpl } from 'react-redux';

import { models, RootModel } from './models';
import { ApiAuth, ApiError } from './models/api';

type FullModel = ExtraModelsFromLoading<RootModel>;

function wrapReducerWithImmer(reducer: Redux.Reducer) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return (state: any, payload: any): any => {
        if (state === undefined) return reducer(state, payload);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        return produce(state, (draft: any) => reducer(draft, payload));
    };
}

const immerPlugin: Plugin<RootModel, FullModel> = {
    onReducer(reducer: Redux.Reducer, _model: string): Redux.Reducer | void {
        return wrapReducerWithImmer(reducer);
    },
};

const loadAuth = () => {
    if (typeof window === 'undefined') {
        return undefined;
    }
    const auth = window.localStorage.getItem('auth');
    return auth ? (JSON.parse(auth) as ApiAuth) : undefined;
};

const loadAuthPlugin: Plugin<RootModel, FullModel> = {
    onStoreCreated(store: Store): Store | void {
        const auth = loadAuth();
        if (auth) {
            store.dispatch({
                type: 'api/setAuth',
                payload: auth,
            });
        }
    },
};

// If a call is made that results in a 401, this will sign the user out.
const badAuthHandler: Redux.Middleware = (store) => (next) => (action) => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const checkError = (err: any) => {
        if (err instanceof ApiError && err.status === 401) {
            store.dispatch({
                type: 'api/setAuth',
                payload: undefined,
            });
            store.dispatch({
                type: 'RESET_ALL',
            });
        } else {
            throw err;
        }
    };

    try {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const p = next(action) as any;
        if (p !== null && typeof p === 'object' && typeof p.then === 'function' && typeof p.catch === 'function') {
            return p.catch(checkError);
        }
        return p;
    } catch (err) {
        checkError(err);
    }
};

export const store = init<RootModel, FullModel>({
    models,
    plugins: [immerPlugin, loadingPlugin(), loadAuthPlugin],
    redux: {
        middlewares: [badAuthHandler],
        rootReducers: {
            RESET_ALL: () => undefined,
        },
    },
});

export type Store = typeof store;
export type Dispatch = RematchDispatch<RootModel>;
export type RootState = RematchRootState<RootModel, FullModel>;

export const useDispatch = () => useDispatchImpl<Dispatch>();
export const useSelector: TypedUseSelectorHook<RootState> = useSelectorImpl;
