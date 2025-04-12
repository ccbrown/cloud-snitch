import type { Metadata } from 'next';
import Link from 'next/link';

import { REVISION } from './revision';

export const metadata: Metadata = {
    title: 'Cookie Policy',
    description: 'Learn how when, how, and why Cloud Snitch uses cookies.',
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4">
            <div className="translucent-snow p-4 rounded-lg">
                <h1>Cookie Policy</h1>
                <div className="text-sm border-b border-platinum pb-4">
                    <span className="uppercase text-english-violet font-semibold">Revision:</span> {REVISION}
                </div>
                <div className="flex flex-col gap-4 mt-4">
                    <p>
                        Our{' '}
                        <Link href="/privacy-policy" className="link">
                            Privacy Policy
                        </Link>{' '}
                        explains our principles when it comes to the collection, processing, and storage of your
                        information. This policy specifically explains how we deploy cookies.
                    </p>
                    <p>We use cookies for very little, so this will be short.</p>
                    <h2>What Are Cookies</h2>
                    <p>
                        Cookies are small pieces of data, stored in text files, that are stored on your computer or
                        other device when websites are loaded in a browser. They are widely used to &quot;remember&quot;
                        you and your preferences, either for a single visit (through a &quot;session cookie&quot;) or
                        for multiple repeat visits (using a &quot;persistent cookie&quot;). They ensure a consistent and
                        efficient experience for visitors, and perform essential functions such as allowing users to
                        register and remain logged in. Cookies may be set by the site that you are visiting (known as
                        &quot;first party cookies&quot;), or by third parties, such as those who serve content or
                        provide advertising or analytics services on the website (&quot;third party cookies&quot;). Both
                        websites and HTML emails may also contain other tracking technologies such as &quot;web
                        beacons&quot; or &quot;pixels.&quot; These are typically small transparent images that provide
                        us with statistics, for similar purposes as cookies. They are often used in conjunction with
                        cookies, though they are not stored on your computer in the same way. As a result, if you
                        disable cookies, web beacons may still load, but their functionality will be restricted.
                    </p>
                    <p>
                        The legal definition of a cookie is broader than the technical definition of a cookie commonly
                        used by engineers. For the purposes of this policy, we use the broader definition which includes
                        tracking pixels, web beacons, local storage, and other technologies that are used to store
                        information on your device.
                    </p>
                    <h2>How We Use Cookies</h2>
                    <p>
                        We use cookies solely to provide our services. We do not use cookies for advertising, analytic,
                        or tracking purposes.
                    </p>
                    <h2>Where We Place Cookies</h2>
                    <p>
                        We set cookies when a user signs into our site and when authenticated users utilize certain
                        features.
                    </p>
                    <h2>Categories of Cookies</h2>
                    <p>There is only one category of cookie that we use: &quot;Essential Cookies&quot;.</p>
                    <p>
                        These cookies are essential for our websites and services to perform basic functions and are
                        necessary for us to operate certain features. These include those required to allow registered
                        users to authenticate and perform account-related functions, store preferences set by users such
                        as account name, language, and location, and ensure our services are operating properly.
                    </p>
                    <h2>Examples</h2>
                    <p>Below are examples of the cookies set by Cloud Snitch, with explanations of their purpose.</p>
                    <ul className="list-disc pl-6">
                        <li>
                            &quot;auth&quot; - This is set when a user signs into the site. It is used to authenticate
                            the user.
                        </li>
                        <li>
                            &quot;team&quot; - This is set whenever a user navigates to a team page. It is used to
                            automatically return the user to the last team they were viewing if they leave and come
                            back.
                        </li>
                    </ul>
                    <h2>Controlling Cookies</h2>
                    <p>
                        Because we only use essential cookies, your only choice if you do not want to use them is to not
                        use our service. It is impossible to use our service without cookies.
                    </p>
                </div>
            </div>
            <div className="text-snow bg-english-violet text-sm p-4 rounded-lg">
                Parts of this agreement have been adapted from Automattic&apos;s{' '}
                <Link
                    href="https://github.com/Automattic/legalmattic/blob/3c220c0969d0ee1731c034ff29fca90460aabf58/Cookie-Policy.md"
                    target="_blank"
                    className="external-snow-link"
                >
                    Cookie Policy
                </Link>{' '}
                under the{' '}
                <Link
                    href="https://creativecommons.org/licenses/by-sa/4.0/"
                    target="_blank"
                    className="external-snow-link"
                >
                    Creative Commons Sharealike license
                </Link>
                .
            </div>
        </div>
    );
};

export default Page;
