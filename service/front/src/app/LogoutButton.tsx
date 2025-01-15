"use client";

import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export default function LogoutButton() {
    const router = useRouter();
    return (
        <Button
            variant='destructive'
            onClick={() => {
                localStorage.removeItem("jwt_token");
                router.push("/login");
            }}
        >
            Logout
        </Button>
    );
}
