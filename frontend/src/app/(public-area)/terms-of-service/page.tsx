import type { Metadata } from 'next';
import Link from 'next/link';

import { REVISION } from './revision';

export const metadata: Metadata = {
    title: 'Terms of Service',
    description: 'These are the terms that govern your access to and use of Cloud Snitch.',
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4">
            <div className="translucent-snow p-4 rounded-lg">
                <h1>Terms of Service</h1>
                <div className="text-sm border-b border-platinum pb-4">
                    <span className="uppercase label">Revision:</span> {REVISION}
                </div>
                <div className="flex flex-col gap-4 mt-4">
                    <p>
                        These Terms govern your access to and use of the software, applications, extensions, and other
                        products and services we provide through or for{' '}
                        <Link href="/" className="link">
                            cloudsnitch.io
                        </Link>
                        .
                    </p>
                    <p>
                        Please read these Terms carefully before accessing or using our Services. By accessing or using
                        any part of our Services, you agree to be bound by all of the Terms and all other operating
                        rules, policies, and procedures that we may publish via the Services from time to time
                        (collectively, the “Agreement”). You also agree that we may automatically change, update, or add
                        on to our Services as stated in the Terms, and the Agreement will apply to any changes.
                    </p>
                    <h2>1. Who&apos;s Who</h2>
                    <p>
                        “You” means any individual or entity using our Services. If you use our Services on behalf of
                        another person or entity, you represent and warrant that you’re authorized to accept the
                        Agreement on that person’s or entity’s behalf, that by using our Services you’re accepting the
                        Agreement on behalf of that person or entity, and that if you, or that person or entity,
                        violates the Agreement, you and that person or entity agree to be responsible to us.
                    </p>
                    <p>
                        Your agreement is with Paragon Cybersecurity LLC (&quot;Paragon Cybersecurity&quot;), a Delaware
                        limited liability company.
                    </p>
                    <h2>2. Your Account</h2>
                    <p>
                        When using our Services requires an account, you agree to provide us with complete and accurate
                        information and to keep the information current so that we can communicate with you about your
                        account. We may need to send you emails about notable updates (like changes to our Terms of
                        Service or Privacy Policy), or to let you know about legal inquiries or complaints we receive
                        about the ways you use our Services so you can make informed choices in response.
                    </p>
                    <p>
                        We may limit your access to our Services until we’re able to verify your account information,
                        like your email address.
                    </p>
                    <p>
                        When you create an account, we consider that to be an inquiry about our products and services,
                        which means that we may also contact you to share more details about what we have to offer
                        (i.e., marketing). Don’t worry — if you aren’t interested, you can opt out of the marketing
                        communications, whether it’s an email, phone call, or text message.
                    </p>
                    <p>
                        You’re solely responsible and liable for all activity under your account. You’re also fully
                        responsible for maintaining the security of your account (which includes keeping your password
                        secure). We’re not liable for any acts or omissions by you, including any damages of any kind
                        incurred as a result of your acts or omissions. If you get fired because of a blog post you
                        write about your boss, that’s on you.
                    </p>
                    <p>
                        Don’t share or misuse your access credentials. And notify us immediately of any unauthorized
                        uses of your account, store, or website, or of any other breach of security. If we believe your
                        account has been compromised, we may suspend or disable it.
                    </p>
                    <p>
                        If you’d like to learn about how we handle the data you provide us, please see our{' '}
                        <Link href="/privacy-policy" className="link">
                            Privacy Policy
                        </Link>
                        .
                    </p>
                    <h2>3. Minimum Age Requirements</h2>
                    <p>
                        Our Services are not directed to children. You’re not allowed to access or use our Services if
                        you’re under the age of 13 (or 16 in Europe). If you register as a user or otherwise use our
                        Services, you represent that you’re at least 13 (or 16 in Europe). You may use our Services only
                        if you can legally form a binding contract with us. In other words, if you’re under 18 years of
                        age (or the legal age of majority where you live), you can only use our Services under the
                        supervision of a parent or legal guardian who agrees to the Agreement.
                    </p>
                    <h2>4. Responsibility of Visitors and Users</h2>
                    <p>
                        We haven’t reviewed, and can’t review, all of the content (like resources names and other
                        materials) posted to or made available through our Services by users or anyone else (”Content”)
                        or on websites that link to, or are linked from, our Services. We’re not responsible for any use
                        or effects of Content or third-party websites.
                    </p>
                    <h2>5. Fees, Payment, and Renewal</h2>
                    <p>
                        Some of our Services are offered for a fee. This section applies to any purchases of Paid
                        Services.
                    </p>
                    <p>
                        By using a Paid Service, you agree to pay the specified fees. Depending on the Paid Service,
                        there may be different kinds of fees, like some that are one-time or recurring. For recurring
                        fees (AKA subscriptions), your subscription begins on your purchase date, and we’ll bill or
                        charge you in the automatically-renewing interval until you cancel, which you can do at any time
                        through the website or by contacting us.
                    </p>
                    <p>
                        To the extent permitted by law, or unless explicitly stated otherwise, all fees do not include
                        applicable federal, provincial, state, local or other governmental sales, value-added, goods and
                        services, harmonized or other taxes, fees, or charges (”Taxes”). You’re responsible for paying
                        all applicable Taxes relating to your use of our Services, your payments, or your purchases. If
                        we’re obligated to pay or collect Taxes on the fees you’ve paid or will pay, you’re responsible
                        for those Taxes, and we may collect payment from you.
                    </p>
                    <p>
                        You must provide accurate and up-to-date payment information. By providing your payment
                        information, you authorize us to store it until you request deletion. If your payment fails, we
                        suspect fraud, or Paid Services are otherwise not paid for or paid for on time (for example, if
                        you contact your bank or credit card company to decline or reverse the charge of fees for Paid
                        Services), we may immediately cancel or revoke your access to Paid Services without notice to
                        you. You authorize us to charge any updated payment information provided by your bank or payment
                        service provider (e.g., new expiration date) or other payment methods provided if we can’t
                        charge your primary payment method.
                    </p>
                    <p>
                        By enrolling in a subscription, you authorize us to automatically charge the then-applicable
                        fees and Taxes for each subsequent subscription period until the subscription is canceled. If
                        you received a discount, used a coupon code, or subscribed during a free trial or promotion,
                        your subscription will automatically renew for the full price of the subscription at the end of
                        the discount period. This means that unless you cancel a subscription, it’ll automatically renew
                        and we’ll charge your payment method(s). You must cancel at least one month before the scheduled
                        end date of any annual subscription and at least 24 hours before the end of any shorter
                        subscription period.
                    </p>
                    <p>
                        You can view your renewal date(s), cancel, or manage subscriptions in your team&apos;s settings
                        or by contacting the support team.
                    </p>
                    <p>
                        We may change our fees at any time in accordance with these Terms and requirements under
                        applicable law. This means that we may change our fees going forward, start charging fees for
                        Services that were previously free, or remove or update features or functionality that were
                        previously included in the fees. If you don’t agree with the changes, you must cancel your Paid
                        Service.
                    </p>
                    <p>
                        We may have a refund policy for some of our Paid Services, and we’ll also provide refunds if
                        required by law. In all other cases, there are no refunds and all payments are final.
                    </p>
                    <h2>6. Feedback</h2>
                    <p>
                        We love hearing from you and are always looking to improve our Services. When you share
                        comments, ideas, or feedback with us, you agree that we’re free to use them without any
                        restriction or compensation to you.
                    </p>
                    <h2>7. General Representation and Warranty</h2>
                    <p>
                        Our mission is to make the web a better place, and our Services are designed to give you control
                        and ownership over your websites. We encourage you to express yourself freely, subject to a few
                        requirements. In particular, you represent and warrant that your use of our Services:
                    </p>
                    <ul className="list-disc pl-6">
                        <li>Will be in strict accordance with the Agreement;</li>
                        <li>
                            Will comply with all applicable laws and regulations (including, without limitation, all
                            applicable laws regarding online conduct and acceptable content, licensing, privacy, data
                            protection, the transmission of technical data exported from the United States or the
                            country in which you reside, the use or provision of financial services, notification and
                            consumer protection, unfair competition, and false advertising);
                        </li>
                        <li>
                            Will not be for any unlawful purposes, to publish illegal content, or in furtherance of
                            illegal activities;
                        </li>
                        <li>
                            Will not infringe or misappropriate the intellectual property rights of Paragon
                            Cybersecurity or any third party;
                        </li>
                        <li>
                            Will not overburden or interfere with our systems or impose an unreasonable or
                            disproportionately large load on our infrastructure, as determined by us in our sole
                            discretion;
                        </li>
                        <li>Will not disclose the personal information of others;</li>
                        <li>Will not be used to send spam or bulk unsolicited messages;</li>
                        <li>Will not interfere with, disrupt, or attack any service or network;</li>
                        <li>
                            Will not be used to create, distribute, or enable material that is, facilitates, or operates
                            in conjunction with, malware, spyware, adware, or other malicious programs or code;
                        </li>
                        <li>
                            Will not involve reverse engineering, decompiling, disassembling, deciphering, or otherwise
                            attempting to derive the source code for the Services or any related technology that is not
                            open source; and
                        </li>
                        <li>
                            Will not involve renting, leasing, loaning, selling, or reselling the Services or related
                            data without our consent.
                        </li>
                    </ul>
                    <h2>8. Changes</h2>
                    <p>
                        We may modify the Terms from time to time, for example, to reflect changes to our Services
                        (e.g., adding new features or benefits to our Services or retiring certain features of certain
                        Services) or for legal, regulatory, or security reasons. If we do this, we’ll provide notice of
                        the changes, such as by posting the amended Terms and updating the “Last Updated” date or, if
                        the changes, in our sole discretion, are material, we may notify you through our Services or
                        other communications. Any changes will apply on a going-forward basis, and, unless we say
                        otherwise, the amended Terms will be effective immediately. By continuing to use our Services
                        after we’ve notified you, you agree to be bound by the new Terms. You have the right to object
                        to any changes at any time by ceasing your use of our Services and canceling any subscription
                        you have.
                    </p>
                    <h2>9. Termination</h2>
                    <p>
                        We may terminate your access to all or any part of our Services at any time, with or without
                        cause or notice, effective immediately, including if we believe, in our sole discretion, that
                        you have violated this Agreement, any service guidelines, or other applicable terms. We have the
                        right (though not the obligation) to (i) refuse or remove any content that, in our reasonable
                        opinion, violates any part of this Agreement, or is in any way harmful or objectionable, (ii)
                        ask you to make adjustments, restrict the resources your website uses, or terminate your access
                        to the Services, if we believe your usage burdens our systems, or (iii) terminate or deny access
                        to and use of any of our Services to any individual or entity for any reason. We will have no
                        obligation to provide a refund of any fees previously paid.
                    </p>
                    <p>
                        You can stop using our Services at any time, or, if you use a Paid Service, you can cancel at
                        any time, subject to the Fees, Payment, and Renewal section of these Terms.
                    </p>
                    <h2>10. Disclaimers</h2>
                    <p>
                        Our Services are provided “as is.” We and our suppliers and licensors hereby disclaim all
                        warranties of any kind, express or implied, to the maximum extent allowed by applicable law,
                        including, without limitation, the warranties of merchantability, fitness for a particular
                        purpose and non-infringement. Neither us, nor our suppliers and licensors, makes any warranty
                        that our Services will be error free or that access thereto will be continuous or uninterrupted.
                        You understand that you download from, or otherwise obtain content or services through, our
                        Services at your own discretion and risk.
                    </p>
                    <h2>11. Jurisdiction and Applicable Law</h2>
                    <p>
                        Except to the extent any applicable law provides otherwise, the Agreement and any access to or
                        use of our Services will be governed by the laws of the state of Delaware, U.S.A., excluding its
                        conflict of law provisions and the application of the United Nations Convention of Contracts for
                        the International Sale of Goods. Nothing in this Agreement affects your rights as a consumer to
                        rely on mandatory provisions in your country of residence.
                    </p>
                    <h2>12. Dispute Resolution</h2>
                    <p>
                        In the event of a dispute arising out of or relating to the Agreement and any access to or use
                        of our Services, parties agree to first attempt to resolve the dispute informally. If the
                        parties are unable to resolve the dispute informally within 30 days, the parties agree to submit
                        the dispute to binding arbitration in accordance with the rules of the American Arbitration
                        Association. The arbitration may be conducted in person, by telephone, or online. Unless
                        otherwise required by law, the arbitration will be conducted in Fulton County, Georgia, U.S.A.
                    </p>
                    <p>
                        In no event shall any dispute arising out of or related to the Agreement or our Services be
                        commenced more than one (1) year after the cause of action arose.
                    </p>
                    <h2>13. Limitation of Liability</h2>
                    <p>
                        In no event will we or our suppliers, partners, or licensors, be liable (including for any
                        third-party products or services purchased or used through our Services) with respect to any
                        subject matter of the Agreement under any contract, negligence, strict liability or other legal
                        or equitable theory for: (i) any special, incidental or consequential damages; (ii) the cost of
                        procurement for substitute products or services; (iii) for interruption of use or loss or
                        corruption of data; or (iv) for any amounts that exceed $50 or the fees paid by you to us under
                        the Agreement during the twelve (12) month period prior to the cause of action, whichever is
                        greater. We shall have no liability for any failure or delay due to matters beyond its
                        reasonable control. The foregoing shall not apply to the extent prohibited by applicable law.
                    </p>
                    <h2>14. Idemnification</h2>
                    <p>
                        You agree to indemnify and hold harmless us, our contractors, and our licensors, and their
                        respective directors, officers, employees, and agents from and against any and all losses,
                        liabilities, demands, damages, costs, claims, and expenses, including attorneys’ fees, arising
                        out of or related to your use of our Services, including but not limited to your violation of
                        the Agreement or any agreement with a provider of third-party services used in connection with
                        the Services or applicable law, Content that you post, and any ecommerce activities conducted
                        through your or another user’s website.
                    </p>
                    <h2>15. Publicity</h2>
                    <p>
                        We may identify you as a customer, using your name and logo in marketing or promotional
                        materials such as customer listings, including on our website. At any time, you may opt out of
                        this or request that we stop using your name and logo by contacting us.
                    </p>
                    <h2>16. Miscellaneous</h2>
                    <p>
                        The Agreement (together with any other terms we provide that apply to any specific Service)
                        constitutes the entire agreement between us and you concerning our Services. If any part of the
                        Agreement is unlawful, void, or unenforceable, that part is severable from the Agreement, and
                        does not affect the validity or enforceability of the rest of the Agreement. A waiver by either
                        party of any term or condition of the Agreement or any breach thereof, in any one instance, will
                        not waive such term or condition or any subsequent breach thereof.
                    </p>
                    <p>
                        We may assign our rights under the Agreement without condition. You may only assign your rights
                        under the Agreement with our prior written consent.
                    </p>
                </div>
            </div>
            <div className="text-snow bg-english-violet text-sm p-4 rounded-lg">
                Parts of this agreement have been adapted from Automattic&apos;s{' '}
                <Link
                    href="https://github.com/Automattic/legalmattic/blob/3c220c0969d0ee1731c034ff29fca90460aabf58/Terms%20of%20Service/WordPress.com/EN-Terms-of-Service.md"
                    target="_blank"
                    className="external-snow-link"
                >
                    Terms of Service
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
