<div align="center">
  <h1>Cloud Snitch <a href="https://github.com/ccbrown/cloud-snitch/actions"><img src="https://img.shields.io/github/actions/workflow/status/ccbrown/cloud-snitch/commit.yaml" alt="GitHub Actions Workflow Status" /></a></h1>

  <p>
    <a href="https://cloudsnitch.io">Get started in minutes at cloudsnitch.io.</a> 🚀
  </p>
</div>

![How well do you really know your cloud?](design/readme/header.png)

Whether you're a developer, security engineer, or just a curious person, Cloud Snitch is guaranteed to teach you something and take your relationship with your cloud to the next level.

Cloud Snitch provides a sleek and intuitive way of exploring your AWS account activity. It's a great addition to any toolbox, regardless of if you're a hobbyist that's just getting started with the cloud or a large enterprise with complex and mature cloud infrastructure.

## Features

![Explore](design/readme/explore.png)

Cloud Snitch gives you a comprehensive look at your AWS account activity in sleek and intuitive map view.

With Cloud Snitch, there's no excuse for not knowing everything your AWS accounts are up to.

![Collaborate](design/readme/collaborate.png)

Invite your teammates to let them explore your AWS activity with you.

Share links to IP address, CIDR network, and AWS principal activity.

Document AWS principals with Markdown notes for your teammates.

![Expose](design/readme/expose.png)

Cloud Snitch provides summaries of activity by AWS region, principal, IP address, and CIDR network.

Errors are highlighted, so you can quickly spot suspicious behavior or bugs in your code.

Take the investigation further with links into to your CloudTrail event history.

## Open Source or SaaS

Cloud Snitch is open sourced under the MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT) so if you're feeling adventurous, you can deploy it directly to your own cloud. Alternatively, you can get up and running in under 5 minutes with an individual or team plan at [cloudsnitch.io](https://cloudsnitch.io).

## Codebase

This repo contains all resources for Cloud Snitch, including the backend, frontend, infrastructure, documentation, and design files. Everything is organized into the following directories:

- [aws](aws): CDK code for deployment to AWS.
- [backend](backend): The backend code and API used by the frontend.
- [design](design): Project files, concepts, screenshots, and other materials not directly required by the codebase.
- [frontend](frontend): The UI for the webapp.

## Contributing

Contributions of all sorts are welcome, but please read the [contribution guidelines](CONTRIBUTING.md) first!

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in the work by you shall be licensed under the MIT license, without any additional terms or conditions.
