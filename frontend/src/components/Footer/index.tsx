import { DefaultFooter } from '@ant-design/pro-layout';
export default () => {
  return (
    <DefaultFooter
      copyright={`${new Date().getFullYear()} xufqing`}
      links={[]}
    />
  );
};
