import { memo } from 'react';

import Logo from './Logo';
import Logout from './Logout';

const Header = ({ title }) => {
    return (
        <header className="bg-blue-800">
            <div className="h-14 sm:h-16 max-w-2xl mx-auto px-4 flex items-center justify-between">
                <div className="flex-1 flex items-center justify-between sm:justify-start">
                    <Logo title={title} />
                    <Logout className="block sm:hidden" />
                </div>

                <Logout className="hidden sm:block" />
            </div>
        </header>
    );
};

export default memo(Header);
