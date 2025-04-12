import { Models } from '@rematch/core';
import { api } from './api';
import { aws } from './aws';
import { reports } from './reports';
import { teams } from './teams';
import { users } from './users';

export interface RootModel extends Models<RootModel> {
    api: typeof api;
    aws: typeof aws;
    reports: typeof reports;
    teams: typeof teams;
    users: typeof users;
}

export const models: RootModel = {
    api,
    aws,
    reports,
    teams,
    users,
};
