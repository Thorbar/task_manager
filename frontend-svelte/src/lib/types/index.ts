export type Task = {
  id: number;
  title: string;
  description: string;
  status: string;
  created_at: string;
  user_id: number;
  username: string;
};


export type User = {
  username: string;
  email: string;
  name: string;
  surname: string;
  created_at: string;
};
