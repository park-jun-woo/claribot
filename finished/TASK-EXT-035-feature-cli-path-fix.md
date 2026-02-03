# TASK-EXT-035: Feature CLI 경로 수정

## 목표
VSCode Extension에서 Feature 추가 시 WSL 터미널에서 올바른 경로로 CLI 명령 실행

## 문제
- `handleCreateFeature`에서 Windows 경로가 WSL 터미널에서 인식되지 않음
- 터미널 `cwd` 옵션이 WSL에서 무시됨

## 수정 내용

### CltEditorProvider.ts
1. Windows 경로를 WSL 경로로 변환하는 함수 추가
2. `handleCreateFeature`에서 터미널에 `cd` 명령으로 디렉토리 이동 후 CLI 실행

```typescript
// Windows 경로를 WSL 경로로 변환
// C:\Users\mail\git\claritask → /mnt/c/Users/mail/git/claritask
function windowsToWslPath(windowsPath: string): string {
  const match = windowsPath.match(/^([A-Za-z]):\\(.*)$/);
  if (match) {
    const drive = match[1].toLowerCase();
    const rest = match[2].replace(/\\/g, '/');
    return `/mnt/${drive}/${rest}`;
  }
  return windowsPath;
}
```

## 테스트
1. VSCode Extension에서 Feature + New 클릭
2. 터미널에서 올바른 경로로 이동 후 clari 명령 실행 확인
