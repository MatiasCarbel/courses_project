export type CourseType = {
  id?: string;
  course_id?: number;
  course_name: string;
  description: string;
  instructor_id: number;
  instructor_name: string;
  category: string;
  requirements: string;
  length: number;
  ImageURL: string;
  CreationTime: string;
  LastUpdated: string;
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
