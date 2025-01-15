"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { ApiResponseType, Department } from "../types";
import { apiUrl } from "../utils";
import { getDepartments } from "./actions";

export default function Page() {
    const [departments, setDepartments] = useState<Department[]>([]);
    const firstNameRef = useRef<HTMLInputElement | null>(null);
    const lastNameRef = useRef<HTMLInputElement | null>(null);
    const emailRef = useRef<HTMLInputElement | null>(null);
    const [department, setDepartment] = useState<number>(0);
    const passwordRef = useRef<HTMLInputElement | null>(null);
    const confirmPasswordRef = useRef<HTMLInputElement | null>(null);
    const router = useRouter();

    useEffect(() => {
        const fetchDepartments = async () => {
            const tmpDepartments = await getDepartments();

            setDepartments(tmpDepartments);
        };
        fetchDepartments();
    }, []);

    async function signup() {
        const firstName = firstNameRef.current?.value;
        const lastName = lastNameRef.current?.value;
        const email = emailRef.current?.value;
        const password = passwordRef.current?.value;
        const confirmPassword = confirmPasswordRef.current?.value;

        const response = await fetch(`${apiUrl}/auth/register`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                first_name: firstName,
                last_name: lastName,
                email,
                department_id: department,
                password,
                confirm_password: confirmPassword,
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
                            <CardTitle className='text-xl'>Sign up</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <form>
                                <div className='grid gap-6'>
                                    <div className='grid gap-6'>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='first_name'>
                                                First Name
                                            </Label>
                                            <Input
                                                id='first_name'
                                                type='text'
                                                placeholder='John'
                                                required
                                                ref={firstNameRef}
                                            />
                                        </div>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='last_name'>
                                                Last Name
                                            </Label>
                                            <Input
                                                id='last_name'
                                                type='text'
                                                placeholder='Doe'
                                                ref={lastNameRef}
                                                required
                                            />
                                        </div>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='department'>
                                                Department
                                            </Label>
                                            <Select
                                                value={
                                                    department == 0
                                                        ? ""
                                                        : department.toString()
                                                }
                                                onValueChange={(value) =>
                                                    setDepartment(
                                                        parseInt(value)
                                                    )
                                                }
                                            >
                                                <SelectTrigger className='w-[180px]'>
                                                    <SelectValue placeholder='Select a department' />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectGroup>
                                                        <SelectLabel>
                                                            Departments
                                                        </SelectLabel>
                                                        {departments.map(
                                                            (department) => (
                                                                <SelectItem
                                                                    value={department.id.toString()}
                                                                    key={
                                                                        department.id
                                                                    }
                                                                >
                                                                    {
                                                                        department.name
                                                                    }
                                                                </SelectItem>
                                                            )
                                                        )}
                                                    </SelectGroup>
                                                </SelectContent>
                                            </Select>
                                        </div>
                                        <div className='grid gap-2'>
                                            <Label htmlFor='email'>Email</Label>
                                            <Input
                                                id='email'
                                                type='email'
                                                placeholder='johndoe@example.com'
                                                ref={emailRef}
                                                required
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
                                                ref={passwordRef}
                                                required
                                            />
                                        </div>
                                        <div className='grid gap-2'>
                                            <div className='flex items-center'>
                                                <Label htmlFor='confirm_password'>
                                                    Confirm Password
                                                </Label>
                                            </div>
                                            <Input
                                                id='confirm_password'
                                                type='password'
                                                ref={confirmPasswordRef}
                                                required
                                            />
                                        </div>
                                        <Button
                                            className='w-full'
                                            type='button'
                                            onClick={signup}
                                        >
                                            Sign up
                                        </Button>
                                    </div>
                                    <div className='text-center text-sm'>
                                        Already have an account?{" "}
                                        <a
                                            href='/signup'
                                            className='underline underline-offset-4'
                                        >
                                            Log in
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
