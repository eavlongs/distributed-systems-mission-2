"use client";
import { Button } from "@/components/ui/button";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { getFiles, getMyProfile } from "./actions";
import DeleteBtn from "./DeleteBtn";
import LogoutButton from "./LogoutButton";
import ProfileDisplay from "./ProfileDisplay";
import SyncNow from "./SyncNow";
import { File, FsType, Profile, RawFile } from "./types";
import UploadFileButton from "./UploadFileButton";
import { apiUrl } from "./utils";

export default function Page() {
    let jwt: string | null = null;
    const [myProfile, setMyProfile] = useState<Profile | null>();
    const [fsToFetch, setFsToFetch] = useState<FsType>(FsType.MAIN);
    const [files, setFiles] = useState<File[] | RawFile[]>([]);
    const [refetch, setRefetch] = useState<boolean>(false);

    const router = useRouter();

    const redirect = () => {
        router.push("/login");
    };

    useEffect(() => {
        const checkLogin = async () => {
            jwt = localStorage.getItem("jwt_token");

            if (jwt === null) {
                redirect();
                return;
            }
            const tmpProfile = await getMyProfile(jwt);
            if (tmpProfile === null) {
                redirect();
                return;
            }
            setMyProfile(tmpProfile);
        };

        checkLogin();
    }, []);

    useEffect(() => {
        const fetchData = async () => {
            jwt = localStorage.getItem("jwt_token");
            const tmpFiles = await getFiles(jwt!, fsToFetch);
            setFiles(tmpFiles);
        };

        fetchData();
    }, [fsToFetch, refetch]);

    return (
        myProfile && (
            <>
                <div className='flex gap-x-4 items-center justify-center mt-2'>
                    <ProfileDisplay profile={myProfile} />
                    <LogoutButton />
                </div>
                <div className='flex justify-center gap-x-4 mt-4'>
                    <SyncNow />
                    <UploadFileButton onUpload={() => setRefetch(!refetch)} />
                </div>
                <div className='flex items-center justify-center mt-4'>
                    {/* <FileUploadForm /> */}
                </div>
                <main className='w-4/5 mx-auto mt-4'>
                    <div className='flex items-center justify-center gap-4'>
                        {Array.of(
                            FsType.MAIN,
                            FsType.BACKUP,
                            FsType.ARCHIVE
                        ).map((tab) => (
                            <Button
                                key={tab}
                                onClick={() => setFsToFetch(tab)}
                                className={`${
                                    fsToFetch === tab
                                        ? "bg-blue-500 text-white hover:bg-blue-600"
                                        : "bg-white text-black hover:bg-gray-200"
                                } px-4 py-2 rounded-md`}
                            >
                                {tab}
                            </Button>
                        ))}
                    </div>
                    {files.length > 0 ? (
                        <div className='rounded-md border'>
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        {fsToFetch == FsType.MAIN && (
                                            <TableHead className='text-center'>
                                                ID
                                            </TableHead>
                                        )}
                                        <TableHead className='text-center'>
                                            Name
                                        </TableHead>
                                        {fsToFetch == FsType.MAIN && (
                                            <TableHead className='w-96 text-center'>
                                                Action
                                            </TableHead>
                                        )}
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {files.map((file) => (
                                        <TableRow key={file.path}>
                                            {fsToFetch == FsType.MAIN && (
                                                <TableCell className='text-center'>
                                                    <span>
                                                        {(file as File).id}
                                                    </span>
                                                </TableCell>
                                            )}
                                            <TableCell className='text-center'>
                                                <a
                                                    href={`${apiUrl}/files/${
                                                        fsToFetch == FsType.MAIN
                                                            ? (file as File).id
                                                            : file.path.replace(
                                                                  "./",
                                                                  ""
                                                              )
                                                    }`}
                                                    target='_blank'
                                                >
                                                    {fsToFetch == FsType.MAIN
                                                        ? (file as File).name
                                                        : file.path.replace(
                                                              `./storage/${myProfile.department_name}/`,
                                                              ""
                                                          )}
                                                </a>
                                            </TableCell>

                                            {fsToFetch == FsType.MAIN && (
                                                <TableCell className='w-96'>
                                                    <div className='flex items-center justify-center'>
                                                        <DeleteBtn
                                                            id={
                                                                (file as File)
                                                                    .id
                                                            }
                                                            onDelete={() =>
                                                                setRefetch(
                                                                    !refetch
                                                                )
                                                            }
                                                        />
                                                    </div>
                                                </TableCell>
                                            )}
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    ) : (
                        <p className='text-center text-2xl font-bold mt-20'>
                            No Files Found
                        </p>
                    )}
                </main>
            </>
        )
    );
}
