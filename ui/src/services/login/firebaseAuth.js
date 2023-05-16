import { validateUser } from "@services/user/validateUser";
import {
    getAuth,
    signInWithPopup,
    GoogleAuthProvider,
    GithubAuthProvider,
} from "firebase/auth";

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

    sessionStorage.setItem("authToken", authToken);

    const userDetails = {
        userId: user?.uid,
        name: user?.displayName,
        email: user?.email,
        avatar: user?.photoURL,
    };

    await validateUser(userDetails);

    return userDetails;
};
