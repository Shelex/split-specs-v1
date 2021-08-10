import { memo } from 'react';

import Header from './Header';
import Footer from './Footer';

const Layout = ({ children }) => {
    return (
        <div className="h-screen w-screen flex flex-col">
            <Header title="Split specs" />

            <main className="flex-1">{children}</main>

            <Footer title="&copy; Shevtsov Oleksandr 2021" />
        </div>
    );
};

export default memo(Layout);
