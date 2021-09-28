import { memo, useEffect } from 'react';
import ApiKeyForm from '../components/CreateApiKeyForm';
import { Link } from 'react-router-dom';

const CreateApiKey = () => {
    return (
        <div className="h-full pt-4 sm:pt-12">
            <ApiKeyForm />
        </div>
    );
};

export default memo(CreateApiKey);
