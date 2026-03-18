const { faker } = require('@faker-js/faker');
const fs = require('fs');
const path = require('path');

// Fixed seed for reproducibility
faker.seed(12345);

const generateUsers = (count = 50) => {
  const users = [];
  for (let i = 1; i <= count; i++) {
    users.push({
      id: i,
      firstName: faker.person.firstName(),
      lastName: faker.person.lastName(),
      email: faker.internet.email(),
    });
  }
  return users;
};

const db = {
  users: generateUsers(),
};

const outputPath = path.join('/data', 'db.json');

try {
  fs.writeFileSync(outputPath, JSON.stringify(db, null, 2));
  console.log(`Successfully generated ${db.users.length} users to ${outputPath}`);
} catch (error) {
  console.error('Error writing db.json:', error);
  process.exit(1);
}
