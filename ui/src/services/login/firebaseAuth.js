import {
    getAuth,
    signInWithPopup,
    GoogleAuthProvider,
    GithubAuthProvider,
} from "firebase/auth";

export const loginWithGoogle = async () => {
    const auth = getAuth();

    signInWithPopup(auth, new GoogleAuthProvider());
};

export const loginWithGithub = async () => {
    const auth = getAuth();

    signInWithPopup(auth, new GithubAuthProvider());
};
