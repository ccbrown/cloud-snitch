import type { Metadata } from 'next';
import Link from 'next/link';

import { REVISION } from './revision';

export const metadata: Metadata = {
    title: 'Privacy Policy',
    description:
        'Your privacy is critically important to us. This is our Privacy Policy, which incorporates and clarifies our principles.',
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4">
            <div className="translucent-snow p-4 rounded-lg">
                <h1>Privacy Policy</h1>
                <div className="text-sm border-b border-platinum pb-4">
                    <span className="uppercase text-english-violet font-semibold">Revision:</span> {REVISION}
                </div>
                <div className="flex flex-col gap-4 mt-4">
                    <p>
                        Your privacy is critically important to us. At Paragon Cybersecurity, we have a few fundamental
                        principles:
                    </p>

                    <ul className="list-disc pl-6">
                        <li>
                            We are thoughtful about the personal information we ask you to provide and the personal
                            information that we collect about you through the operation of our services.
                        </li>
                        <li>We store personal information for only as long as we have a reason to keep it.</li>
                        <li>We help protect you from overreaching government demands for your personal information.</li>
                        <li>
                            We aim for full transparency on how we gather, use, and share your personal information.
                        </li>
                    </ul>

                    <p>
                        This is our Privacy Policy, which incorporates and clarifies these principles. It applies to
                        information that we collect about you when you use our services.
                    </p>
                    <p>
                        Below we explain how we collect, use, and share information about you, along with the choices
                        that you have with respect to that information.
                    </p>
                    <h2>Information We Collect</h2>
                    <p>
                        We only collect information about you if we have a reason to do so — for example, to provide our
                        Services, to communicate with you, or to make our Services better.
                    </p>
                    <p>
                        We collect this information from three sources: if and when you provide information to us,
                        automatically through operating our Services, and from outside sources. Let&apos;s go over the
                        information that we collect.
                    </p>
                    <h3>Information You Provide to Us</h3>
                    <p>
                        It&apos;s probably no surprise that we collect information that you provide to us directly. Here
                        are some examples:
                    </p>
                    <ul className="list-disc pl-6">
                        <li>
                            <strong>Basic account information:</strong>We ask for basic information from you in order to
                            set up your account. For example, we require individuals who sign up for an account to
                            provide an email address and optionally a password — and that&apos;s it. You may provide us
                            with more information in order to use certain features of our services, but we don&apos;t
                            require that information to create an account.
                        </li>
                        <li>
                            <strong>Payment information:</strong> When signing up for a subscription, we&apos;ll collect
                            information required to process payments. This will include a name and address for billing
                            purposes. Financial details such as credit card information are processed by a third party,{' '}
                            <Link href="http://stripe.com" target="_blank" className="external-link">
                                Stripe
                            </Link>
                            . When provided, those details will be securely handled by Stripe and will at no point be
                            visible to our servers. We also keep a record of the transactions you&apos;ve made.
                        </li>
                        <li>
                            <strong>Team names and data:</strong> When creating a team, you&apos;ll provide it with a
                            name, which doesn&apos;t have to correspond to any real individual or organization name.
                            You&apos;ll also likely want to configure an integration with your cloud provider, which
                            will cause our service to collect data from your cloud environment. This data is used to
                            provide you with the service, is not shared with any third parties, and can be deleted at
                            any time. We won&apos;t even look at it unless it&apos;s strictly necessary to support you.
                        </li>
                        <li>
                            <strong>Communications with us (hi there!):</strong> You may also provide us with
                            information when you respond to surveys, communicate with us about a support question, or
                            post a question via public forums. When you communicate with us via form, email, phone, or
                            otherwise, we store a copy of our communications (including any call recordings as permitted
                            by applicable law).
                        </li>
                    </ul>
                    <h3>Information We Collect Automatically</h3>
                    <p>We also collect some information automatically:</p>
                    <ul className="list-disc pl-6">
                        <li>
                            <strong>Log information:</strong> Like most online service providers, we collect information
                            that web browsers, mobile devices, and servers typically make available, including the
                            browser type, IP address, unique device identifiers, language preference, referring site,
                            the date and time of access, operating system, and mobile network information.
                        </li>
                        <li>
                            <strong>Transactional information:</strong> When you make a purchase through our Services,
                            we collect information about the transaction, such as product details, purchase price, and
                            the date and location of the transaction.
                        </li>
                        <li>
                            <strong>Usage information:</strong> We collect information about your usage of our Services.
                            For example, we collect information about the actions that site administrators and users
                            perform on a site — in other words, who did what and when. We also collect information about
                            what happens when you use our Services (e.g., page views) along with information about your
                            device (e.g., screen size, name of cellular network, and mobile device manufacturer). We use
                            this information to, for example, provide our Services to you, get insights on how people
                            use our Services so we can make our Services better, and understand and make predictions
                            about user retention.
                        </li>
                        <li>
                            <strong>Location information:</strong> We may determine the approximate location of your
                            device from your IP address. We collect and use this information to, for example, calculate
                            how many people visit our Services from certain geographic regions.
                        </li>
                        <li>
                            <strong>Information from cookies:</strong> A cookie is a string of information that a
                            website stores on a visitor&apos;s computer, and that the visitor&apos;s browser provides to
                            the website each time the visitor returns. We use cookies only for the purposes of providing
                            our services once a user has logged in. We do not use cookies to track users across the web
                            or for any other purpose.
                        </li>
                    </ul>
                    <h3>Information We Collect from Other Sources</h3>
                    <p>
                        Third-party services may also give us information, like mailing addresses for individuals who
                        are not yet our users (but we hope will be!). We use this information for marketing purposes
                        like postcards and other mailers advertising our Services.
                    </p>
                    <h2>How and Why We Use Information</h2>
                    <h3>Purposes for Using Information</h3>
                    <p>We use information about you for the purposes listed below:</p>
                    <ul className="list-disc pl-6">
                        <li>
                            <strong>To provide our Services.</strong> For example, to set up and maintain your account,
                            provide customer service, process payments and orders, and verify user information.
                        </li>
                        <li>
                            <strong>To ensure quality, maintain safety, and improve our Services.</strong> For example,
                            by providing automatic upgrades and new versions of our Services. Or, for example, by
                            monitoring and analyzing how users interact with our Services so we can create new features
                            that we think our users will enjoy and that will help them create and manage websites more
                            efficiently or make our Services easier to use.
                        </li>
                        <li>
                            <strong>
                                To market our Services and measure, gauge, and improve the effectiveness of our
                                marketing.
                            </strong>{' '}
                            For example, by targeting our marketing messages to groups of our users (like those who have
                            a particular plan with us or have been users for a certain length of time), advertising our
                            Services, analyzing the results of our marketing campaigns (like how many people purchased a
                            paid plan after receiving a marketing message), and understanding and forecasting user
                            retention.
                        </li>
                        <li>
                            <strong>To protect our Services, our users, and the public.</strong> For example, by
                            detecting security incidents; detecting and protecting against malicious, deceptive,
                            fraudulent, or illegal activity; fighting spam; complying with our legal obligations; and
                            protecting the rights and property of Paragon Cybersecurity and others, which may result in
                            us, for example, declining a transaction or terminating Services.
                        </li>
                        <li>
                            <strong>To fix problems with our Services.</strong> For example, by monitoring, debugging,
                            repairing, and preventing issues.
                        </li>
                        <li>
                            <strong>To customize the user experience.</strong> For example, to personalize your
                            experience by serving you relevant notifications and advertisements for our Services.
                        </li>
                        <li>
                            <strong>To communicate with you.</strong> For example, by emailing you to ask for your
                            feedback, share tips for getting the most out of our products, or keep you up to date on
                            Paragon Cybersecurity; texting you to verify your payment; or calling you to share offers
                            and promotions that we think will be of interest to you. If you don&apos;t want to hear from
                            us, you can opt out of marketing communications at any time. (If you opt out, we&apos;ll
                            still send you important updates relating to your account.)
                        </li>
                    </ul>
                    <h3>Legal Bases for Processing Information</h3>
                    <p>
                        A note here for those in the European Union about our legal grounds for processing information
                        about you under EU data protection laws, which is that our use of your information is based on
                        the grounds that:
                    </p>
                    <ul className="list-disc pl-6">
                        <li>
                            The use is necessary in order to fulfill our commitments to you under the applicable terms
                            of service or other agreements with you or is necessary to administer your account — for
                            example, in order to device or charge you for a paid plan; or
                        </li>
                        <li>The use is necessary for compliance with a legal obligation; or</li>
                        <li>
                            The use is necessary in order to protect your vital interests or those of another person; or
                        </li>
                        <li>
                            We have a legitimate interest in using your information — for example, to provide and update
                            our Services; to improve our Services so that we can offer you an even better user
                            experience; to safeguard our Services; to communicate with you; to measure, gauge, and
                            improve the effectiveness of our advertising; and to understand our user retention and
                            attrition; to monitor and prevent any problems with our Services; and to personalize your
                            experience; or
                        </li>
                        <li>
                            You have given us your consent — for example before we place certain cookies on your device
                            and access and analyze them later on as described in our{' '}
                            <Link href="/cookie-policy" className="link">
                                Cookie Policy
                            </Link>
                            .
                        </li>
                    </ul>
                    <h2>Sharing Information</h2>
                    <h3>How We Share Information</h3>
                    <p>
                        We share information about you in limited circumstances, and with appropriate safeguards on your
                        privacy. These are spelled out below.
                    </p>
                    <ul className="list-disc pl-6">
                        <li>
                            <strong>Subsidiaries and independent contractors:</strong> We may disclose information about
                            you to our subsidiaries and independent contractors who need the information to help us
                            provide our Services or process the information on our behalf. We require our subsidiaries
                            and independent contractors to follow this Privacy Policy for any personal information that
                            we share with them.
                        </li>
                        <li>
                            <strong>Third-party vendors:</strong> We may share information about you with third-party
                            vendors who need the information in order to provide their services to us, or to provide
                            their services to you or your site. This includes vendors that help us provide our Services
                            to you (like Stripe, postal and email delivery services that help us stay in touch with you,
                            customer chat and email support services that help us communicate with you); those that
                            assist us with our marketing efforts (e.g., by providing tools for identifying a specific
                            marketing target group or improving our marketing campaigns, and by placing ads to market
                            our services); those that help us understand and enhance our Services (like analytics
                            providers); those that make tools to help us run our operations (like programs that help us
                            with task management, scheduling, word processing, email and other communications, and
                            collaboration among our teams); and other third-party tools that help us manage operations.
                            We require vendors to agree to privacy commitments in order to share information with them.
                        </li>
                        <li>
                            <strong>Legal and regulatory requirements:</strong> We may disclose information about you in
                            response to a subpoena, court order, or other governmental request.
                        </li>
                        <li>
                            <strong>To protect rights, property, and others:</strong> We may disclose information about
                            you when we believe in good faith that disclosure is reasonably necessary to protect the
                            property or rights of Paragon Cybersecurity, third parties, or the public at large. For
                            example, if we have a good faith belief that there is an imminent danger of death or serious
                            physical injury, we may disclose information related to the emergency without delay.
                        </li>
                        <li>
                            <strong>Business transfers:</strong> In connection with any merger, sale of company assets,
                            or acquisition of all or a portion of our business by another company, or in the unlikely
                            event that Paragon Cybersecurity goes out of business or enters bankruptcy, user information
                            would likely be one of the assets that is transferred or acquired by a third party. If any
                            of these events were to happen, this Privacy Policy would continue to apply to your
                            information and the party receiving your information may continue to use your information,
                            but only consistent with this Privacy Policy.
                        </li>
                        <li>
                            <strong>With your consent:</strong> We may share and disclose information with your consent
                            or at your direction. For example, we may share your information with third parties when you
                            authorize us to do so.
                        </li>
                        <li>
                            <strong>Aggregated or de-identified information:</strong> We may share information that has
                            been aggregated or de-identified, so that it can no longer reasonably be used to identify
                            you. For instance, we may publish aggregate statistics about the use of our Services, or
                            share a hashed version of your email address to facilitate customized ad campaigns on other
                            platforms.
                        </li>
                        <li>
                            <strong>Published support requests:</strong> If you send us a request for assistance (for
                            example, via a support email or one of our other feedback mechanisms), we reserve the right
                            to publish that request in order to clarify or respond to your request, or to help us
                            support other users.
                        </li>
                    </ul>
                    <h3>Information Shared With Collaborators</h3>
                    <p>
                        Our services provide various tools for sharing information with your designated team members.
                        For example, team settings and billing information may be shared with other team administrators.
                    </p>
                    <p>
                        Please be mindful when deciding who to invite to your team and what information you share with
                        them.
                    </p>
                    <h2>How Long We Keep Information</h2>
                    <p>
                        We generally discard information about you when it&apos;s no longer needed for the purposes for
                        which we collect and use it — described in the section above on How and Why We Use Information —
                        and we&apos;re not legally required to keep it.
                    </p>

                    <p>
                        For example, we keep web server logs that record information about a visitor to Cloud Snitch,
                        like the visitor&apos;s IP address, browser type, and operating system, for approximately 30
                        days. We retain the logs for this period of time in order to, among other things, analyze
                        traffic and investigate issues if something goes wrong on one of our websites.
                    </p>

                    <p>
                        Information collected from cloud providers is retained based on the retention policy of your
                        team&apos;s subscription.
                    </p>
                    <h2>Security</h2>
                    <p>
                        While no online service is 100% secure, we work very hard to protect information about you
                        against unauthorized access, use, alteration, or destruction, and take reasonable measures to do
                        so. We monitor our Services for potential vulnerabilities and attacks.
                    </p>
                    <h2>Choices</h2>
                    <p>You have several choices available when it comes to information about you:</p>
                    <ul className="list-disc pl-6">
                        <li>
                            <strong>Limit the information that you provide:</strong> If you have an account with us, you
                            can choose not to provide the optional account information, profile information, and
                            transaction and billing information. Please keep in mind that if you do not provide this
                            information, certain features of our Services may not be accessible.
                        </li>
                        <li>
                            <strong>Close your account:</strong> While we&apos;d be very sad to see you go, you can
                            close your account if you no longer want to use our Services. Please keep in mind that we
                            may continue to retain your information after closing your account, as described in How Long
                            We Keep Information above — for example, when that information is reasonably needed to
                            comply with (or demonstrate our compliance with) legal obligations such as law enforcement
                            requests, or reasonably needed for our legitimate business interests.
                        </li>
                    </ul>
                    <h2>Your Rights</h2>
                    <p>
                        If you are located in certain parts of the world, including some US states and countries that
                        fall under the scope of the European General Data Protection Regulation (aka the
                        &quot;GDPR&quot;), you may have certain rights regarding your personal information, like the
                        right to request access to or deletion of your data.
                    </p>
                    <h3>European General Data Protection Regulation (GDPR)</h3>
                    <p>
                        If you are located in a country that falls under the scope of the GDPR, data protection laws
                        give you certain rights with respect to your personal data, subject to any exemptions provided
                        by the law, including the rights to:
                    </p>
                    <ul className="list-disc pl-6">
                        <li>Request access to your personal data;</li>
                        <li>Request correction or deletion of your personal data;</li>
                        <li>Object to our use and processing of your personal data;</li>
                        <li>Request that we limit our use and processing of your personal data; and</li>
                        <li>Request portability of your personal data.</li>
                    </ul>
                    <p>You also have the right to make a complaint to a government supervisory authority.</p>
                    <h3>Contacting Us About These Rights</h3>
                    <p>
                        You can usually access, correct, or delete your personal data using your account settings and
                        tools that we offer, but if you aren&apos;t able to or you&apos;d like to contact us about one
                        of the other rights, see{' '}
                        <Link href="/contact" className="link">
                            Contact Us
                        </Link>{' '}
                        to, well, find out how to contact us.
                    </p>
                    <p>
                        When you contact us about one of your rights under this section, we&apos;ll need to verify that
                        you are the right person before we disclose or delete anything. For example, if you are a user,
                        we will need you to contact us from the email address associated with your account. You can also
                        designate an authorized agent to make a request on your behalf by giving us written
                        authorization. We may still require you to verify your identity with us.
                    </p>
                    <h3>Appeals Process for Rights Requests Denials</h3>
                    <p>
                        In some circumstances we may deny your request to exercise one of these rights. For example, if
                        we cannot verify that you are the account owner we may deny your request to access the personal
                        information associated with your account. As another example, if we are legally required to
                        maintain a copy of your personal information we may deny your request to delete your personal
                        information.
                    </p>
                    <p>
                        In the event that we deny your request, we will communicate this fact to you in writing. You may
                        appeal our decision by responding in writing to our denial email and stating that you would like
                        to appeal. All appeals will be reviewed by an internal expert who was not involved in your
                        original request. In the event that your appeal is also denied this information will be
                        communicated to you in writing. Please note that the appeal process does not apply to job
                        applicants.
                    </p>
                    <p>
                        If your appeal is denied, in some US states you may refer the denied appeal to the state
                        attorney general if you believe the denial is in conflict with your legal rights. The process
                        for how to do this will be communicated to you in writing at the same time we send you our
                        decision about your appeal.
                    </p>
                    <h2>Privacy Policy Changes</h2>
                    <p>
                        Although most changes are likely to be minor, we may change its Privacy Policy from time to
                        time. We encourages visitors to frequently check this page for any changes to its Privacy
                        Policy. If we make changes, we will notify you by revising the change log below, and, in some
                        cases, we may provide additional notice sending you a notification through email or your
                        dashboard. Your further use of the Services after a change to our Privacy Policy will be subject
                        to the updated policy.
                    </p>
                </div>
            </div>
            <div className="text-snow bg-english-violet text-sm p-4 rounded-lg">
                Parts of this agreement have been adapted from Automattic&apos;s{' '}
                <Link
                    href="https://github.com/Automattic/legalmattic/blob/3c220c0969d0ee1731c034ff29fca90460aabf58/Privacy-Policy.md"
                    target="_blank"
                    className="external-snow-link"
                >
                    Privacy Policy
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
