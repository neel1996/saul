import {
    getAuth,
    signInWithPopup,
    GoogleAuthProvider,
    GithubAuthProvider,
} from "firebase/auth";

import { login } from "./userLogin";

export const loginWithGoogle = async () => {
    const auth = getAuth();

    return await signInWithPopup(auth, new GoogleAuthProvider()).then(
        authenticatedUser
    );
};

export const loginWithGithub = async () => {
    const auth = getAuth();

    return await signInWithPopup(auth, new GithubAuthProvider()).then(
        authenticatedUser
    );
};

const authenticatedUser = async (response) => {
    const { user } = response;
    const authToken = user.accessToken;

    if (!authToken) {
        throw new Error("Invalid credentials");
    }

    sessionStorage.setItem("authToken", authToken);

    const userDetails = {
        userId: user?.uid,
        name: user?.displayName,
        email: user?.email,
        avatar: user?.photoURL,
    };

    await login(userDetails);

    return userDetails;
};
