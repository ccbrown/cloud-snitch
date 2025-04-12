import type { Metadata } from 'next';

import { MascotBox } from '@/components';
import { ContactForm } from './ContactForm';

export const metadata: Metadata = {
    title: 'Contact Us',
    description:
        'Submit a support request, ask a question about Cloud Snitch, or just say hi! We love hearing from you.',
};

const Page = () => {
    return (
        <div>
            <MascotBox>
                <div className="flex flex-col gap-4">
                    <h1>Contact Us</h1>
                    <ContactForm />
                </div>
            </MascotBox>
        </div>
    );
};

export default Page;
