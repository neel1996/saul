import { loginWithGoogle, loginWithGithub } from "@services/login";
import "@testing-library/jest-dom";
import { act, fireEvent, render, screen } from "@testing-library/react";
import React from "react";
import { redirect } from "react-router-dom";

import Login from "./Login";

jest.mock("@services/login");
jest.mock("react-router-dom", () => ({
    ...jest.requireActual("react-router-dom"),
    redirect: jest.fn(),
}));

describe("Login", () => {
    afterEach(() => {
        jest.clearAllMocks();
        jest.resetAllMocks();
    });

    it("should render login component correctly", async () => {
        const { container } = await render(<Login />);

        expect(await screen.findByTestId("login")).toBeInTheDocument();

        expect(container).toMatchSnapshot();
    });

    it("should login with google on click", async () => {
        loginWithGoogle.mockResolvedValue({});

        await render(<Login />);

        const loginWithGoogleButton = await screen.findByTestId(
            "login-with-google"
        );

        // eslint-disable-next-line testing-library/no-unnecessary-act
        await act(async () => {
            await fireEvent.click(loginWithGoogleButton);
        });

        expect(loginWithGoogle).toHaveBeenCalledTimes(1);
        expect(redirect).toHaveBeenCalledTimes(1);
        expect(redirect).toHaveBeenCalledWith("/");
    });

    it("should login with github on click", async () => {
        loginWithGithub.mockResolvedValue({});

        await render(<Login />);

        const loginWithGithubButton = await screen.findByTestId(
            "login-with-github"
        );

        // eslint-disable-next-line testing-library/no-unnecessary-act
        await act(async () => {
            await fireEvent.click(loginWithGithubButton);
        });

        expect(loginWithGithub).toHaveBeenCalledTimes(1);
        expect(redirect).toHaveBeenCalledTimes(1);
        expect(redirect).toHaveBeenCalledWith("/");
    });
});
