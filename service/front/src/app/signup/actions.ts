import { ApiResponseType, Department } from "../types";
import { apiUrl } from "../utils";

export async function getDepartments() {
    try {
        console.log(`${apiUrl}/departments`);
        const response = await fetch(`${apiUrl}/departments`);
        const json: ApiResponseType<Department[]> = await response.json();

        if (response.status !== 200 || !json.success) {
            throw new Error(json.message);
        }

        return json.data!;
    } catch (e: any) {
        console.log(e.message);
        return [];
    }
}
