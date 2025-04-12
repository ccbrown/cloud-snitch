import { createModel } from '@rematch/core';

import { RootModel } from '.';
import { apiConfiguration } from './api';

import {
    BeginUserEmailAuthenticationInput,
    BeginUserRegistrationInput,
    CompleteUserPasskeyRegistrationInput,
    UpdateUserInput,
    User,
    UserApi,
    UserPasskey,
} from '@/generated/api';

interface UsersState {
    currentUserId: string | null;
    users: Record<string, User>;
    passkeys: Record<string, Record<string, UserPasskey>>;
}

export const users = createModel<RootModel>()({
    state: {
        currentUserId: null,
        users: {},
        passkeys: {},
    } as UsersState,
    reducers: {
        put(state, user: User) {
            state.users[user.id] = user;
        },
        setCurrentUserId(state, id: string) {
            state.currentUserId = id;
        },
        putPasskey(state, passkey: UserPasskey) {
            if (!state.passkeys[passkey.userId]) {
                state.passkeys[passkey.userId] = {};
            }
            state.passkeys[passkey.userId][passkey.id] = passkey;
        },
        setUserPasskeys(state, userId: string, passkeys: UserPasskey[]) {
            const m: Record<string, UserPasskey> = {};
            passkeys.forEach((p) => {
                m[p.id] = p;
            });
            state.passkeys[userId] = m;
        },
        removePasskey(state, passkeyId: string) {
            Object.values(state.passkeys).forEach((ps) => {
                delete ps[passkeyId];
            });
        },
    },
    effects: (dispatch) => ({
        async fetchCurrent(_payload: void, state) {
            if (!state.api.auth) {
                return;
            }
            const api = new UserApi(apiConfiguration(state.api));
            const resp = await api.getUser({
                userId: 'self',
            });
            dispatch.users.put(resp);
            dispatch.users.setCurrentUserId(resp.id);
        },
        async fetch(id: string, state) {
            const api = new UserApi(apiConfiguration(state.api));
            const resp = await api.getUser({
                userId: id,
            });
            dispatch.users.put(resp);
        },
        async updateUser(input: { id: string; update: UpdateUserInput }, state) {
            const api = new UserApi(apiConfiguration(state.api));
            const resp = await api.updateUser({
                userId: input.id,
                updateUserInput: input.update,
            });
            dispatch.users.put(resp);
        },
        async beginRegistration(payload: BeginUserRegistrationInput, state) {
            const api = new UserApi(apiConfiguration(state.api));
            await api.beginUserRegistration({
                beginUserRegistrationInput: payload,
            });
        },
        async beginEmailAuthentication(payload: BeginUserEmailAuthenticationInput, state) {
            const api = new UserApi(apiConfiguration(state.api));
            await api.beginUserEmailAuthentication({
                beginUserEmailAuthenticationInput: payload,
            });
        },
        async completeRegistration(
            input: {
                token: string;
            },
            state,
        ) {
            const api = new UserApi(apiConfiguration(state.api));
            const resp = await api.completeUserRegistration({
                completeUserRegistrationInput: {
                    token: input.token,
                },
            });
            dispatch.users.put(resp.user);
            dispatch.users.setCurrentUserId(resp.user.id);
            dispatch.api.setAuth({
                token: resp.accessToken,
            });
        },
        async fetchPasskeys(userId: string, state) {
            const api = new UserApi(apiConfiguration(state.api));
            const resp = await api.getUserPasskeys({
                userId,
            });
            dispatch.users.setUserPasskeys(userId, resp);
        },
        async beginPasskeyAuthentication(_payload: void, state) {
            const api = new UserApi(apiConfiguration(state.api));
            return await api.beginUserPasskeyAuthentication({
                body: {},
            });
        },
        async beginPasskeyRegistration(userId: string, state) {
            const api = new UserApi(apiConfiguration(state.api));
            return await api.beginUserPasskeyRegistration({
                body: {},
                userId,
            });
        },
        async completePasskeyRegistration(
            input: {
                userId: string;
                registration: CompleteUserPasskeyRegistrationInput;
            },
            state,
        ) {
            const api = new UserApi(apiConfiguration(state.api));
            const resp = await api.completeUserPasskeyRegistration({
                userId: input.userId,
                completeUserPasskeyRegistrationInput: input.registration,
            });
            dispatch.users.putPasskey(resp);
        },
        async deletePasskey(id: string, state) {
            const api = new UserApi(apiConfiguration(state.api));
            await api.deleteUserPasskeyById({
                passkeyId: id,
            });
            dispatch.users.removePasskey(id);
        },
    }),
});
