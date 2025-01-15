import { Profile } from "./types";

export default function ProfileDisplay({ profile }: { profile: Profile }) {
    return (
        <div className='flex flex-row gap-x-4 justify-between items-center p-4 border border-gray-300 rounded-lg'>
            <p className='text-lg font-semibold'>
                {profile.first_name} {profile.last_name}
            </p>
            <p className='text-sm text-gray-500'>{profile.email}</p>

            <p className='text-sm text-gray-500'>
                Department: {profile.department_name}
            </p>
        </div>
    );
}
