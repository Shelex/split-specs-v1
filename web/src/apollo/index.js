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
        if (networkError.statusCode === 401) {
            unsetToken();
        }
        console.error(`[Network error]: ${networkError}`);
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
