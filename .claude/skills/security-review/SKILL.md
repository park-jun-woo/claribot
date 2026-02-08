---
name: security-review
description: 서버/인프라 보안 검토. bastion, AWS SG, nginx, SSH, 포트, 인증서 등 전체 보안 점검
argument-hint: [target - bastion, aws, all]
allowed-tools: Bash(ssh *), Bash(aws *), Bash(curl *), Bash(nmap *), Bash(openssl *), Read, Grep, Glob
---

보안전문가로서 인프라 보안을 검토한다.
대상: $ARGUMENTS (미지정 시 all)

## 검토 대상별 체크리스트

### bastion (또는 all)

SSH로 bastion 서버에 접속하여 점검:
- SSH key: `ssh -i /mnt/c/Users/mail/git/bastion/bastion-key.pem ubuntu@43.201.60.23`

**SSH 보안:**
- `PasswordAuthentication no` 확인
- `PermitRootLogin no` 확인
- `GatewayPorts no` 확인
- SSH 포트 변경 여부
- authorized_keys 검토 (불필요 키 없는지)

**방화벽 (UFW):**
- `sudo ufw status verbose` - 기본 정책 deny인지
- 불필요 포트 열려있지 않은지
- outbound 규칙 최소화 확인

**네트워크:**
- `sudo ss -tlnp` - 리스닝 포트 확인
- localhost만 바인딩해야 할 서비스가 0.0.0.0에 열려있지 않은지
- tunnel 포트(19847)가 127.0.0.1에만 바인딩 확인

**nginx:**
- `server_tokens off` 확인
- 시크릿 헤더 검증 동작 확인
- 불필요한 sites-enabled 없는지
- 버전 정보 노출 여부 (curl -sI)

**OS:**
- `cat /etc/os-release` - OS 버전
- `sudo apt list --upgradable` - 보안 패치 대기 중인지
- `sudo lastlog` - 최근 로그인 기록
- `sudo cat /var/log/auth.log | tail -30` - 인증 로그

### aws (또는 all)

AWS CLI로 점검 (credentials 필요):

**Security Group:**
- `aws ec2 describe-security-groups --group-ids sg-0491cffbbff2104da` - 인바운드 규칙 최소화 확인
- 22번 포트가 특정 IP만 허용하는지
- 불필요한 0.0.0.0/0 규칙 없는지

**CloudFront:**
- 배포 상태, TLS 버전, 캐시 정책 확인
- 시크릿 헤더 설정 확인

**ACM:**
- 인증서 만료일 확인
- DNS 검증 상태 확인

**IAM:**
- 사용 중인 IAM 사용자/역할 권한 과다 여부

### 외부 접근 테스트

bastion IP 직접 접근 차단 확인:
```
curl -sI http://43.201.60.23:9847/
curl -sf https://clari.parkjunwoo.com/api/health
```

## 출력 형식

검토 결과를 표로 정리:

| 항목 | 상태 | 비고 |
|------|------|------|
| ... | 양호/주의/위험 | 설명 |

마지막에 **조치 필요 사항**이 있으면 즉시 수정 가능한 건 바로 수정하고, 사용자 확인이 필요한 건 제안한다.
