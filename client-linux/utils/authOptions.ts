import { type NextAuthOptions, type User as NextAuthUser } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import jwt from "jsonwebtoken";
import { JWT } from "next-auth/jwt";

const authOptions: NextAuthOptions = {
  providers: [
    CredentialsProvider({
      name: "credentials",
      credentials: {
        email: { label: "Email", type: "text" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        if (!credentials || !credentials.email || !credentials.password) {
          return null;
        }

        try {
          const response = await fetch(`${process.env.BACKEND_URL}/login`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              email: credentials.email,
              password: credentials.password,
            }),
          });
          const responseData = await response.json();
          if (response.ok) {
            return {
              id: responseData.user.id,
              name: responseData.user.name,
              email: responseData.user.email,
              avatar: responseData.user.avatar,
            } satisfies NextAuthUser;
          }

          return null;
        } catch (error) {
          console.log(error);
          return null;
        }
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user, session, trigger }) {
      if (user) {
        token.id = user.id;
        token.name = user.name;
        token.email = user.email;
        token.avatar = user.avatar;
      }

      if (trigger === "update" && session) {
        token.name = session.name ?? token.name;
        token.avatar = session.avatar ?? token.avatar;
      }

      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        session.user.id = token.id;
        session.user.name = token.name;
        session.user.email = token.email;
        session.user.avatar = token.avatar;
      }
      return session;
    },
  },
  session: {
    strategy: "jwt",
  },
  jwt: {
    async encode({ token, secret }) {
      return jwt.sign(token as object, secret, { algorithm: "HS256" });
    },
    async decode({ token, secret }) {
      return jwt.verify(token!, secret) as JWT;
    },
  },
  pages: {
    signIn: "/auth?mode=login",
  },
  secret: process.env.NEXTAUTH_SECRET,
  cookies: {
    sessionToken: {
      name: "next-auth.session-token",
      options: {
        httpOnly: false,
        sameSite: "lax",
        path: "/",
        secure: process.env.NODE_ENV === "production",
      },
    },
  },
};

export default authOptions;
