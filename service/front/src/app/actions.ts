import { ApiResponseType, File, FsType, Profile, RawFile } from "./types";
import { apiUrl, archiveUrl, backupUrl } from "./utils";

export async function deleteFile(token: string, id: string) {
    try {
        const response = await fetch(`${apiUrl}/files/${id}`, {
            headers: {
                Authorization: "Bearer " + token,
            },
            method: "DELETE",
        });

        const json: ApiResponseType = await response.json();

        if (response.status !== 200 || !json.success) {
            throw new Error(json.message);
        }

        return true;
    } catch (e: any) {
        console.log(e.message);
        return false;
    }
}

export async function getMyProfile(token: string): Promise<Profile | null> {
    try {
        const response = await fetch(`${apiUrl}/auth/whoami`, {
            headers: {
                Authorization: "Bearer " + token,
            },
        });

        const json: ApiResponseType<Profile> = await response.json();

        if (response.status !== 200 || !json.success) {
            throw new Error(json.message);
        }

        return json.data!;
    } catch (e: any) {
        console.log(e.message);
        return null;
    }
}

export async function getFiles(
    token: string,
    fsToFetch: FsType
): Promise<File[] | RawFile[]> {
    try {
        let apiUrlToFetch = "";
        switch (fsToFetch) {
            case FsType.MAIN:
                apiUrlToFetch = apiUrl as string;
                break;
            case FsType.BACKUP:
                apiUrlToFetch = backupUrl as string;
                break;
            case FsType.ARCHIVE:
                apiUrlToFetch = archiveUrl as string;
                break;
        }

        const response = await fetch(`${apiUrlToFetch}/files/${fsToFetch}`, {
            headers: {
                Authorization: "Bearer " + token,
            },
        });

        const json: ApiResponseType<{ files: File[] | RawFile[] }> =
            await response.json();

        if (response.status !== 200 || !json.success) {
            throw new Error(json.message);
        }

        return json.data!.files;
    } catch (e: any) {
        console.log(e.message);
        return [];
    }
}
