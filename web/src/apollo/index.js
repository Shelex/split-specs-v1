import {
    ApolloClient,
    createHttpLink,
    InMemoryCache,
    makeVar,
    from
} from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { setContext } from '@apollo/client/link/context';

const token = () => localStorage.getItem('token');
const unsetToken = () => localStorage.removeItem('token');

const httpLink = createHttpLink({
    uri: 'https://split-specs.appspot.com/query'
});

const errorLink = onError(({ graphQLErrors, networkError, response }) => {
    if (networkError) {
        console.error(`[Network error]: ${networkError}`);
    }

    if (graphQLErrors) {
        graphQLErrors.map(({ message, locations, path }) => {
            if (!path.includes('nextSpec')) {
                window.location.href = '/';
            }
            if (message.includes('access denied')) {
                unsetToken();
            } else {
                return alert(
                    `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
                );
            }
        });
    }
});

const authLink = setContext((_, { headers }) => ({
    headers: {
        ...headers,
        Authorization: token() ?? ''
    }
}));

const cache = new InMemoryCache({
    typePolicies: {
        Query: {
            fields: {
                isLoggedIn: {
                    read() {
                        return isLoggedInVar();
                    }
                }
            }
        }
    }
});

const client = new ApolloClient({
    link: from([errorLink, authLink, httpLink]),
    cache
});

export const isLoggedInVar = makeVar(Boolean(token()));

export default client;
