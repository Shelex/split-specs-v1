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

// TODO:
//   changePassword(input: ChangePasswordInput!): String!
//   shareProject(email: String!, projectName: String!): String!
//   deleteSession(sessionId: String!): String!
//   deleteProject(projectName: String!): String!
