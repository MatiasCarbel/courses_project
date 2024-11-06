import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const data = await request.json();
  const { Email, PasswordHash, Username, FirstName, LastName, UserType } = data;

  const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";

  const usersReq = await fetch(`${baseUrl}/user/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: Email,
      password_hash: PasswordHash,
      username: Username,
      name: FirstName,
      last_name: LastName,
      usertype: UserType,
    }),
  });

  const usersJson = await usersReq.json();
  if (usersJson.error)
    return NextResponse.json({ message: usersJson.error }, { status: 401 });

  return NextResponse.json({ message: "Created." }, { status: 200 });
}
