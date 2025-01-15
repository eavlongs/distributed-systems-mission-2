export type Profile = {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    department_id: string;
    department_name: string;
};

export type ApiResponseType<T = any> = {
    success: boolean;
    error?: any;
    data?: T;
    message: string;
};

export enum FsType {
    MAIN = "main",
    BACKUP = "backup",
    ARCHIVE = "archive",
}

export type File = {
    id: string;
    name: string;
    path: string;
    departmentID: string;
    is_deleted: boolean;
    created_at: Date;
    updated_at: Date;
};

export type RawFile = Pick<File, "path">;

export type Department = {
    id: number;
    name: string;
};
