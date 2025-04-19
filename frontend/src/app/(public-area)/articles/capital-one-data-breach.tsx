const article = {
    title: 'Capital One Data Breach – A Cautionary Tale',
    description:
        'Capital One lost hundreds of millions after being notified by a third party of an intruder that had been lurking in their AWS account for four months.',
    date: new Date(Date.parse('2025-04-15T23:10:00-04:00')),
    content: (
        <>
            <p>
                In July of 2019, Capital One Financial Corporation announced that they had recently become aware of a
                data breach compromising the personal data of over 100 million customers in the United States, including
                140,000 Social Security numbers and 80,000 bank account numbers. Capital One estimated that the direct
                financial impact of the incident would be between $100 million and $150 million for notifications,
                credit monitoring, technology, and legal support. Furthermore, trust in the company had plummeted, and
                within two weeks of the incident, Capital One&apos;s share price had dropped by 14%.
            </p>
            <p>
                The impact was substantial, but it might have been much less severe if not for the fact that{' '}
                <strong>the breach had occurred almost four months before it was discovered</strong>. Capital One had
                been moving its operations to the cloud, and had failed to adequately manage the risk of doing so. This
                mistake would later cost them another $80 million in fines when all was said and done.
            </p>
            <p>
                So what should they have done to manage their risk better? Let&apos;s break down the incident in detail
                so we can learn from their mistakes.
            </p>
            <h2>Foothold</h2>
            <p>
                On March 22, 2019, the attacker discovered that they could trick Capital One EC2 servers running in AWS
                into executing arbitrary HTTP requests, enabling a type of attack known as Server Side Request Forgery
                (SSRF). They were able to do this as a result of a misconfigured web application firewall (WAF) known as
                ModSecurity. The WAF allowed the attacker to send requests to the AWS metadata service, which among
                other things provides the EC2 servers with temporary credentials for the purpose of invoking AWS APIs.
            </p>
            <p>
                At this point, the attacker could impersonate the EC2 servers and perform any action that the associated
                IAM role allowed.
            </p>
            <h2>Lateral Movement</h2>
            <p>
                Once the attacker gained their foothold, they began to enumerate everything they had access to. Usually,{' '}
                <strong>this is loud</strong> as attackers typically just have to try things and see what works. If you
                were to look at this role in Cloud Snitch, you would see it attempting (and hopefully failing) to do a
                lot of things that you would not expect a WAF to do.
            </p>
            <p>
                In Capital One&apos;s case, the permissions attached to the role allowed the attacker to read and
                decrypt sensitive data from Amazon S3 buckets. The attacker was able to exfiltrate all of this data via
                the WAF&apos;s credentials, which again should have sounded alarm bells.
            </p>
            <h2>Discovery</h2>
            <p>
                Capital One had completely failed to detect the intrusion. They only became aware of it when a third
                party notified them via their Responsible Disclosure Program on July 17, 2019.
            </p>
            <p>
                Capital One was completely oblivious to the fact that there was an intruder in their cloud for almost
                four months. The presence of the intruder would have made a lot of noise in their CloudTrail logs, but
                without the proper tools, that log activity went completely unnoticed. This is why Cloud Snitch was
                created.
            </p>
            <p>
                <strong>
                    With Cloud Snitch, you can easily confirm that your AWS resources are doing what you expect them to
                    do and nothing more.
                </strong>
            </p>
        </>
    ),
    relatedLinks: [
        {
            title: 'Capital One Announces Data Security Incident',
            url: 'https://investor.capitalone.com/news-releases/news-release-details/capital-one-announces-data-security-incident/',
        },
        {
            title: 'Information on the Capital One cyber incident',
            url: 'https://www.capitalone.com/digital/facts2019/',
        },
        {
            title: 'A Systematic Analysis of the Capital One Data Breach: Critical Lessons Learned',
            url: 'https://dl.acm.org/doi/10.1145/3546068',
        },
        {
            title: 'Capital One to pay $80 million fine after data breach',
            url: 'https://www.reuters.com/article/business/capital-one-to-pay-80-million-fine-after-data-breach-idUSKCN2522D8/',
        },
        {
            title: 'Capital One Data Breach — 2019',
            url: 'https://medium.com/nerd-for-tech/capital-one-data-breach-2019-f85a259eaa60',
        },
        {
            title: 'Lessons from the Capital One Breach on Cloud Security',
            url: 'https://www.darktrace.com/blog/back-to-square-one-the-capital-one-breach-proved-we-must-rethink-cloud-security',
        },
    ],
};

export default article;
