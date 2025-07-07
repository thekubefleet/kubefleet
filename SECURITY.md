# Security Policy

## Supported Versions

Use this section to tell people about which versions of your project are
currently being supported with security updates.

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of KubeFleet seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Reporting Process

1. **Do not create a public GitHub issue** for the vulnerability
2. **Email us** at [security@kubefleet.io](mailto:security@kubefleet.io) with the details
3. **Include the following information**:
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact
   - Suggested fix (if any)
   - Your contact information

### What to Expect

- **Acknowledgement**: You will receive an acknowledgment within 48 hours
- **Assessment**: We will assess the reported vulnerability within 7 days
- **Updates**: We will keep you informed of our progress
- **Resolution**: We will work to resolve the issue and release a fix

### Responsible Disclosure

We follow responsible disclosure practices:

- **Timeline**: We aim to fix critical vulnerabilities within 30 days
- **Credit**: Security researchers will be credited in our security advisories
- **Coordination**: We will coordinate with you on the disclosure timeline

### Security Best Practices

When using KubeFleet, please follow these security best practices:

1. **Keep updated**: Always use the latest stable version
2. **Network security**: Use TLS for gRPC communication
3. **RBAC**: Follow the principle of least privilege for Kubernetes permissions
4. **Monitoring**: Monitor the agent and dashboard for unusual activity
5. **Access control**: Restrict access to the dashboard in production

### Security Features

KubeFleet includes several security features:

- **gRPC TLS**: Support for encrypted communication
- **RBAC integration**: Proper Kubernetes role-based access control
- **Input validation**: All inputs are validated and sanitized
- **Error handling**: Secure error handling without information disclosure

### Known Issues

Currently, there are no known security vulnerabilities in KubeFleet.

### Security Updates

Security updates will be released as patch versions (e.g., 1.0.1, 1.0.2) and will be clearly marked in the release notes.

### Contact Information

For security-related issues, please contact:
- **Email**: [security@kubefleet.io](mailto:security@kubefleet.io)
- **PGP Key**: [Available upon request]

Thank you for helping keep KubeFleet secure! 🔒 