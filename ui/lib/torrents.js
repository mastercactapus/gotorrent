
type ActiveTorrent = {
  ID: string;
  Name: string;
  Length: number;
  Chunks: number;
  BytesCompleted: number;
  Seeding: boolean;
  Downspeed: number;
  ETA: number;
}
