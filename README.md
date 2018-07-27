# mitmengine

The goal of this project is to allow for accurate detection of HTTPS interception and robust TLS fingerprinting. 
This project is based off of [The Security Impact of HTTPS Interception](https://zakird.com/papers/https_interception.pdf), and started as a port to Go of [their processing scripts and fingerprints](https://github.com/zakird/tlsfingerprints).

## Signature and Fingerprints
In this project, fingerprints map to concrete instantiations of an object, while signatures can represent multiple objects. We use this convention because a fingerprint is usually an inherent property of an object, while a signature can be chosen. In the same way, an actual client request seen by a server would have a fingerprint, while the software generating the request can choose it's own signature (e.g., by choosing which cipher suites it supports).

### Client Request Fingerprint
A client request fingerprint is derived from a client request to a server, and contains both TLS and HTTP features.

### Client Request Signature
A client request signature represents a set of possible request fingerprints. The aim is to make each signature specific enough that it can uniquely identify a piece of software.

### User Agent Fingerprint
A user agent fingerprint is derived from the raw user agent in a client request. It consists of the browser name, browser version, OS platform, OS name, OS version, device type, and a 'quirks' field that contains additional flags. For example, if 'GSA/' (Google Search Appliance) is present in the raw user agent, the quirk 'gsa' is appended to the user agent quirks.

### User Agent Signature
A user agent signature represent a set of user agent fingerprints. It consists of the browser name, browser version range (min-max), OS platform, OS name, OS version range (min-max), device type, and a quirk signature. Ideally, a user agent should identify a specific browser or software version.

### Mitm Info
The MITM info field contains additional information about known man-in-the-middle software. This field contains the software name, classification (e.g., antivirus, proxy), and an additional maximum security grade which can be set to take into account factors such as whether or not the software validates certificates.

### Record
Records that are loaded into the MITM detection database consist of a user agent signature, a request signature, and a mitm info field (which is empty for signatures of legitimate software).

## MITM Detection Methodology
Incoming client requests are first matched against a database of known-good
browser signatures to get all records with a user agent signature that matches
the request's user agent fingerprint. If any of the returned records contain a
request signature that matches the incoming client request's signature, then
the request is marked as 'possible'. Otherwise, the request is marked as
'unlikely' or 'impossible' and the MITM engine matches the request fingerprint
against a database of known MITM software signatures to see if there is a
match.

### False positives
If a signature is inaccurate or outdated for a given piece of client software,
it is possible that the signature will falsely flag a connection as being
intercepted.

### False negatives
If a proxy closely mimics the request of the client, then we may not expect to
detect a mismatch. If the browser signatures are overly broad, we will also
fail to detect interception.

## Testing
To test, run ```make test``` and to see code coverage, run ```make cover```.

## Godoc
Run ```make godoc``` or ```PKG=<sub-package> make godoc``` to generate godoc for mitmengine or any of the sub-packages.

## API
The intended entrypoint to the mitmengine package is through the `Processor.Check` function, which takes a user agent and client request fingerprint, and returns a mitm detection report. Additional API functions may be added in the future to allow for adding new signatures to a running process, for example.
