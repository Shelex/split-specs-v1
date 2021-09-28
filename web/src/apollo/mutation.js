import { gql } from '@apollo/client';

export const SIGN_UP = gql`
    mutation register($email: String!, $password: String!) {
        register(input: { email: $email, password: $password })
    }
`;

export const SIGN_IN = gql`
    mutation login($email: String!, $password: String!) {
        login(input: { email: $email, password: $password })
    }
`;

export const DELETE_SESSION = gql`
    mutation deleteSession($sessionId: String!) {
        deleteSession(sessionId: $sessionId)
    }
`;

export const DELETE_PROJECT = gql`
    mutation deleteProject($projectName: String!) {
        deleteProject(projectName: $projectName)
    }
`;

export const CREATE_SESSION = gql`
    mutation addSession($session: SessionInput!) {
        addSession(session: $session) {
            sessionId
            projectName
        }
    }
`;

export const CREATE_API_KEY = gql`
    mutation addApiKey($name: String!, $expireAt: Int!) {
        addApiKey(name: $name, expireAt: $expireAt)
    }
`;

export const DELETE_API_KEY = gql`
    mutation deleteApiKey($keyId: String!) {
        deleteApiKey(keyId: $keyId)
    }
`;

// TODO:
//   changePassword(input: ChangePasswordInput!): String!
//   shareProject(email: String!, projectName: String!): String!
