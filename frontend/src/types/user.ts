export type UserResponse = {
  id: string;
  access_token: string;
  expires_at: number;
};

export type MeResponse = {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  language: string;
};
