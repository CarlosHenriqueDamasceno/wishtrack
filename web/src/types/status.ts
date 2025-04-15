export enum Status {
  UNSAVED,
  SAVED,
  RATED,
}

export const statusName: Record<Status, string> = {
  [Status.UNSAVED]: 'unsaved',
  [Status.SAVED]: 'saved',
  [Status.RATED]: 'rated',
}
