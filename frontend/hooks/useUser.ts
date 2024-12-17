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
      const response = await fetch("/api/auth/user");
      const data = await response.json();

      if (response.ok && data.user) {
        setUser(data.user);
        setIsAdmin(data.user.admin === true);
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
    await fetch("/api/auth/logout");
    setUser(null);
    setIsAdmin(false);
  };

  const refreshUser = () => {
    checkUser();
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
