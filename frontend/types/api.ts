export type ListResult<T> = {
  data: T[],
  size: number
}

export type TournamentRow = {
  ExternalID: string,
  Name: string,
  StartDate: string,
  EndDate: string,
  Tier: string,
  Location: string,
  Country: string
}
