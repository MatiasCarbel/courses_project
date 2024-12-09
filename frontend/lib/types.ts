export type CourseType = {
  id?: string;
  title: string;
  description: string;
  instructor: string;
  category: string;
  duration: number;
  available_seats: number;
  image_url: string;
  is_subscribed: boolean;
};

export type UserType = {
  id: number;
  email: string;
  username: string;
  name: string;
  last_name: string;
  usertype: boolean;
  password_hash: string;
  CreationTime: string;
  LastUpdated: string;
};

export type CommentType = {
  comment_id: number;
  course_id: number;
  user_id: number;
  user_name: string;
  comment: string;
  CreationTime: string;
  LastUpdated: string;
};
