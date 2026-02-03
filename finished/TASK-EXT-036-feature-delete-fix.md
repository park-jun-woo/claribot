# TASK-EXT-036: Feature Delete 버튼 수정

## 목표
VSCode Extension에서 Feature 삭제 시 관련된 task_edges도 함께 삭제

## 문제
- `deleteFeature`에서 tasks 삭제 전에 task_edges를 삭제하지 않아 외래키 제약 조건 위반

## 수정 내용

### database.ts
1. `deleteFeature` 메서드에서 task_edges 삭제 추가

```typescript
deleteFeature(featureId: number): void {
  // Get all task IDs for this feature
  const tasks = this.queryAll<{ id: number }>(
    'SELECT id FROM tasks WHERE feature_id = ?',
    [featureId]
  );
  const taskIds = tasks.map(t => t.id);

  // Delete related task_edges first
  if (taskIds.length > 0) {
    const placeholders = taskIds.map(() => '?').join(',');
    this.run(
      `DELETE FROM task_edges WHERE from_task_id IN (${placeholders}) OR to_task_id IN (${placeholders})`,
      [...taskIds, ...taskIds]
    );
  }

  // Delete feature_edges
  this.run('DELETE FROM feature_edges WHERE from_feature_id = ? OR to_feature_id = ?', [featureId, featureId]);
  // Delete tasks
  this.run('DELETE FROM tasks WHERE feature_id = ?', [featureId]);
  // Delete feature
  this.run('DELETE FROM features WHERE id = ?', [featureId]);
  this.save();
}
```

## 테스트
1. Feature 선택 후 Delete 버튼 클릭
2. Feature와 관련 Tasks가 정상적으로 삭제되는지 확인
