"use client";

import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { deleteFile } from "./actions";

export default function DeleteBtn({
    id,
    onDelete,
}: {
    id: string;
    onDelete: () => void;
}) {
    let token = localStorage.getItem("jwt_token") as string;
    return (
        <Button
            variant='destructive'
            size='icon'
            onClick={async () => {
                const success = await deleteFile(token, id);

                if (success) {
                    onDelete();
                } else {
                    alert("Error deleting file");
                }
            }}
        >
            <Trash2 />
        </Button>
    );
}
