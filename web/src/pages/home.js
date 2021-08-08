import { useReactiveVar } from '@apollo/client';
import { memo, useEffect } from 'react';
import { Redirect } from 'react-router-dom';

import { isLoggedInVar } from '../apollo';
import SignInForm from '../components/SignInForm';

const Home = () => {
    const isLoggedIn = useReactiveVar(isLoggedInVar);

    useEffect(() => {
        document.title = 'Split Specs';
    }, []);

    return (
        <section>
            {isLoggedIn ? <Redirect to="/projects" /> : <SignInForm />}
        </section>
    );
};

export default memo(Home);
