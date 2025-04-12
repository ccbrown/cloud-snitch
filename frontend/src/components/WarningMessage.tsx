import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';

export const WarningMessage = ({ children }: { children: React.ReactNode }) => (
    <p className="w-full my-2 rounded-md p-3 text-dark-purple bg-yellow-100 sm:text-sm sm:leading-6 [&_a]:underline">
        <span className="inline-flex items-baseline">
            <ExclamationTriangleIcon className="h-6 w-6 shrink-0 self-center mr-2" />
            <span>{children}</span>
        </span>
    </p>
);
