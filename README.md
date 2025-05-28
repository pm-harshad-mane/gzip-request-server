# gzip-request-server

test.html has code that makes a plain/text post request to Golang server (main.go)
test_zip.html code makes the gzip post request sent in contetType = "plain/text" to Golang server (main_zip.go)
Diff of main.go and main_zip.go displays the minimal changes required to support this.

Linkedin post: https://www.linkedin.com/feed/update/urn:li:activity:7331425419407085568/

Prebid documentation: https://docs.prebid.org/dev-docs/bidder-adaptor.html#compression-support-for-outgoing-requests 

Prebid JS Commit: https://github.com/prebid/Prebid.js/commit/77b620da035d9365d0091beafc268b78a215a9fd#diff-44300c00ba8a077145a3d52158f46e60a6fecfef87aa2ccbc1fe70ee148e8dc7

Code change in PubMatic bidder:
<img width="532" alt="Screenshot 2025-05-28 at 10 27 40 AM" src="https://github.com/user-attachments/assets/69ace66e-10c9-4442-8786-b4456d23db2e" />

Code diff of required server side change:
<img width="1728" alt="Screenshot 2025-05-27 at 7 21 05 AM" src="https://github.com/user-attachments/assets/108043b1-d447-4611-9961-14417d4c991b" />
