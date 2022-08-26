import { DatabaseInfo } from './dbinfo';
export interface OpenedDatabaseConnection {
    conid: string;
    database: string;
    structure: DatabaseInfo;
    analysedTime?: number;
    serverVersion?: any;
    disconnected?: boolean;
    status?: {
        name: string;
        message?: string;
        counter: number;
    };
}

export * from './dbinfo.d';
export * from './appdefs';