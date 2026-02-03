import type { MessageFromWebview } from './types';

declare function acquireVsCodeApi(): {
  postMessage(message: any): void;
  getState(): any;
  setState(state: any): void;
};

const vscodeApi = acquireVsCodeApi();

export const vscode = {
  postMessage: (message: MessageFromWebview) => vscodeApi.postMessage(message),
};

export function postMessage(message: MessageFromWebview): void {
  vscodeApi.postMessage(message);
}

export function getState<T>(): T | undefined {
  return vscodeApi.getState() as T | undefined;
}

export function setState<T>(state: T): void {
  vscodeApi.setState(state);
}
