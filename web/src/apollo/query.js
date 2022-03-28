import { gql } from '@apollo/client';

export const GET_PROJECTS = gql`
    query projects {
        projects
    }
`;

export const GET_PROJECT = gql`
    query project($name: String!, $pagination: Pagination) {
        project(name: $name, pagination: $pagination) {
            projectName
            totalSessions
            sessions {
                id
                start
                end
                backlog {
                    file
                    estimatedDuration
                    assignedTo
                    start
                    end
                }
            }
        }
    }
`;

export const GET_SESSION = gql`
    query session($id: String!) {
        session(sessionId: $id) {
            id
            start
            end
            backlog {
                file
                estimatedDuration
                assignedTo
                start
                end
                passed
            }
        }
    }
`;

export const NEXT_SPEC = gql`
    query nextSpec($sessionId: String!, $options: NextOptions) {
        nextSpec(sessionId: $sessionId, options: $options)
    }
`;

export const API_KEYS = gql`
    query {
        getApiKeys {
            id
            name
            expireAt
        }
    }
`;
