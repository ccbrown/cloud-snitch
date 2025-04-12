import type { Metadata } from 'next';

import { MascotBox } from '@/components';
import { RegistrationForm } from './RegistrationForm';

export const metadata: Metadata = {
    title: 'Register',
    description: 'Create a Cloud Snitch account and take your relationship with your cloud to the next level.',
};

const Page = () => {
    return (
        <div>
            <MascotBox>
                <div className="flex flex-col gap-4">
                    <h1>Register</h1>
                    <RegistrationForm />
                </div>
            </MascotBox>
        </div>
    );
};

export default Page;
