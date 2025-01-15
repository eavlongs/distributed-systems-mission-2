"use client";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useRouter } from "next/navigation";
import { useRef } from "react";
import { ApiResponseType } from "../types";
import { apiUrl } from "../utils";

export default function Page() {
    const emailRef = useRef<HTMLInputElement>(null);
    const passwordRef = useRef<HTMLInputElement>(null);
    const router = useRouter();

    async function login() {
        const email = emailRef.current?.value;
        const password = passwordRef.current?.value;

        const response = await fetch(`${apiUrl}/auth/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email,
                password,
            }),
        });

        const json: ApiResponseType<{
            token: string;
        }> = await response.json();

        if (response.ok) {
            localStorage.setItem("jwt_token", json.data!.token);
            router.push("/");
        } else {
            alert(json.message);
        }
    }

    return (
        <div className='flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10'>
            <div className='flex w-full max-w-sm flex-col gap-6'>
                <div className='flex flex-col gap-6'>
                    <Card>
                        <CardHeader className='text-center'>
                            <CardTitle className='text-xl'>Log in</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <form>
                                <div className='grid gap-6'>
                                    <div className='grid gap-6'>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='email'>Email</Label>
                                            <Input
                                                id='email'
                                                type='email'
                                                placeholder='johndoe@example.com'
                                                required
                                                ref={emailRef}
                                            />
                                        </div>
                                        <div className='grid gap-2'>
                                            <div className='flex items-center'>
                                                <Label htmlFor='password'>
                                                    Password
                                                </Label>
                                            </div>
                                            <Input
                                                id='password'
                                                type='password'
                                                required
                                                ref={passwordRef}
                                            />
                                        </div>
                                        <Button
                                            type='button'
                                            className='w-full'
                                            onClick={login}
                                        >
                                            Login
                                        </Button>
                                    </div>
                                    <div className='text-center text-sm'>
                                        Don&apos;t have an account?{" "}
                                        <a
                                            href='/signup'
                                            className='underline underline-offset-4'
                                        >
                                            Sign up
                                        </a>
                                    </div>
                                </div>
                            </form>
                        </CardContent>
                    </Card>
                    <div className='text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 [&_a]:hover:text-primary  '>
                        By clicking continue, you agree to our{" "}
                        <a href='#'>Terms of Service</a> and{" "}
                        <a href='#'>Privacy Policy</a>.
                    </div>
                </div>
            </div>
        </div>
    );
}
