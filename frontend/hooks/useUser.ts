import { useState, useEffect } from "react";
import { jwtDecode, JwtPayload } from "jwt-decode";

interface UserJwtPayload extends JwtPayload {
  id: string;
  username: string;
  email: string;
  admin: boolean;
}

export const useUser = () => {
  const [user, setUser] = useState<UserJwtPayload | null>(null);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const checkUser = async () => {
    try {
      const response = await fetch("/api/auth/user", {
        cache: "no-store",
        credentials: "include",
        headers: {
          "Cache-Control": "no-cache",
          Pragma: "no-cache",
        },
      });

      const data = await response.json();

      if (data.user) {
        if (data.user.username && data.user.user_id) {
          setUser(data.user);
          setIsAdmin(data.user.admin === true);
        } else {
          console.error("Incomplete user data received:", data.user);
          setUser(null);
          setIsAdmin(false);
        }
      } else {
        setUser(null);
        setIsAdmin(false);
      }
    } catch (error) {
      console.error("Error checking user:", error);
      setUser(null);
      setIsAdmin(false);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    checkUser();
  }, []);

  const logout = async () => {
    try {
      await fetch("/api/auth/logout");
      setUser(null);
      setIsAdmin(false);
    } catch (error) {
      console.error("Error during logout:", error);
    }
  };

  const refreshUser = async () => {
    setIsLoading(true);
    await checkUser();
  };

  return {
    user,
    isAdmin,
    isLoading,
    isAuthed: !!user,
    logout,
    refreshUser,
  };
};
