import { Button } from "@/components/ui/button";
import { Label } from "@radix-ui/react-label";
import { ChangeEvent, useRef } from "react";
import { apiUrl } from "./utils";
import { ApiResponseType } from "./types";

export default function UploadFileButton({
    onUpload,
}: {
    onUpload: () => void;
}) {
    let token = localStorage.getItem("jwt_token");
    const fileInputRef = useRef<HTMLInputElement>(null);

    async function uploadFile(e: ChangeEvent<HTMLInputElement>) {
        const file = e.target.files![0];

        if (!file) {
            return;
        }

        const formData = new FormData();
        formData.append("file", file);

        const response = await fetch(`${apiUrl}/files/upload`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
            method: "POST",
            body: formData,
        });

        const json: ApiResponseType = await response.json();
        if (!response.ok) {
            alert(json.message);
        }

        onUpload();
        e.target.files = null;
        e.target.value = "";
    }

    return (
        <form>
            <input
                type='file'
                hidden
                id='file_upload'
                name='file_upload'
                onChange={uploadFile}
            />

            <Button type='button'>
                <Label htmlFor='file_upload'>Upload File</Label>
            </Button>
        </form>
    );
}
