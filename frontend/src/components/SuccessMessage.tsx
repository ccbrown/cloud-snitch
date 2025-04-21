import { CheckCircleIcon } from '@heroicons/react/24/outline';

export const SuccessMessage = ({ children }: { children: React.ReactNode }) => (
    <p className="w-full rounded-md p-3 text-dark-purple bg-emerald/50 sm:text-sm sm:leading-6 [&_a]:underline">
        <span className="inline-flex items-baseline">
            <CheckCircleIcon className="h-6 w-6 shrink-0 self-center mr-2" />
            <span>{children}</span>
        </span>
    </p>
);
