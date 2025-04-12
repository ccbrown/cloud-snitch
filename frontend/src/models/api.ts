import { createModel } from '@rematch/core';
import { RootModel } from '.';

import { Configuration, ModelError, ResponseContext, UserApi, UserCredentials } from '@/generated/api';

export interface ApiAuth {
    token: string;
}

interface ApiState {
    auth?: ApiAuth;
}

const getBasePath = () => process.env.NEXT_PUBLIC_API_URL || `${window.location.protocol}//${window.location.host}/api`;

export class ApiError extends Error {
    status: number;

    constructor(message: string, status: number) {
        super(message);
        this.status = status;
    }
}

const middleware = [
    {
        post: async (context: ResponseContext) => {
            if (context.response.status !== 200) {
                const body = (await context.response.json()) as ModelError;
                throw new ApiError(body.message, context.response.status);
            }
            return context.response;
        },
    },
];

export const apiConfiguration = (state?: ApiState): Configuration => {
    return new Configuration({
        basePath: getBasePath(),
        headers: state?.auth && { Authorization: `token ${state.auth.token}` },
        middleware,
    });
};

export const api = createModel<RootModel>()({
    state: {
        auth: undefined,
    } as ApiState,
    reducers: {
        setAuth(state, auth: ApiAuth | undefined) {
            state.auth = auth;
            if (auth) {
                window.localStorage.setItem('auth', JSON.stringify(auth));
            } else {
                window.localStorage.removeItem('auth');
            }
        },
    },
    effects: (dispatch) => ({
        async signIn(payload: UserCredentials) {
            const userApi = new UserApi(
                new Configuration({
                    basePath: getBasePath(),
                    middleware,
                }),
            );
            const resp = await userApi.authenticate({
                userCredentials: payload,
            });
            dispatch.api.setAuth({
                token: resp.token,
            });
        },
        async signOut(_payload: void, state) {
            if (state.api.auth) {
                const api = new UserApi(apiConfiguration(state.api));
                await api.signOut({});
                dispatch.api.setAuth(undefined);
                await dispatch({ type: 'RESET_ALL' });
            }
        },
    }),
});
